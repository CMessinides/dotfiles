#! /usr/bin/env python3
import argparse
from enum import Enum
from functools import total_ordering
import os
from pathlib import Path
import subprocess
import sys
from typing import Any


PACKAGE_ROOT = Path("./packages")
HOME = Path.home()


class FileStatus(Enum):
    UNINSTALLED = "uninstalled"
    BROKEN = "broken"
    CONFLICTING = "conflicting"
    INSTALLED = "installed"


class PackageStatus(Enum):
    EMPTY = "empty"
    UNINSTALLED = "uninstalled"
    PARTIALLY_INSTALLED = "partially installed"
    BROKEN = "broken"
    INSTALLED = "installed"


def get_file_status(pkg: "Package", file: Path) -> FileStatus:
    target = HOME / file
    real_target = target.resolve()

    if not target.exists():
        # target either doesn't exist, or there's a broken symlink somewhere along its path
        is_broken = target.is_symlink() or (
            target != real_target and not real_target.exists()
        )

        if is_broken:
            return FileStatus.BROKEN
        else:
            return FileStatus.UNINSTALLED

    if real_target != (pkg.path / file).resolve():
        # target exists, but it isn't a link to our package
        return FileStatus.CONFLICTING
    else:
        # target does exist and links to our package
        return FileStatus.INSTALLED


class PackageState:
    def __init__(self, pkg: "Package", files: list[Path]) -> None:
        self.pkg = pkg
        self.files: dict[str, FileStatus] = {}

        if pkg.exists():
            for file in files:
                self.files[str(file)] = get_file_status(pkg, file)

    def is_empty(self) -> bool:
        return self.status() == PackageStatus.EMPTY

    def is_installed(self) -> bool:
        return self.status() == PackageStatus.INSTALLED

    def is_uninstalled(self) -> bool:
        return self.status() == PackageStatus.UNINSTALLED

    def is_broken(self) -> bool:
        return self.status() == PackageStatus.BROKEN

    def can_install(self) -> bool:
        return not self.is_broken()

    def status(self) -> PackageStatus:
        statuses = set(self.files.values())

        if len(statuses) == 0:
            return PackageStatus.EMPTY
        elif statuses == {FileStatus.INSTALLED}:
            return PackageStatus.INSTALLED
        elif statuses == {FileStatus.UNINSTALLED}:
            return PackageStatus.UNINSTALLED
        elif statuses == {FileStatus.INSTALLED, FileStatus.UNINSTALLED}:
            return PackageStatus.PARTIALLY_INSTALLED
        else:
            return PackageStatus.BROKEN

    def __str__(self) -> str:
        lines: list[str] = [f"{self.pkg.name} ({self.status().value})"]

        for file, status in self.files.items():
            lines.append(f"  {file} ({status.value})")

        return "\n".join(lines)

    def __repr__(self) -> str:
        return f"<PackageState status={self.status()}>"


@total_ordering
class Package:
    def __init__(self, name: str) -> None:
        self.name = name
        self.path = PACKAGE_ROOT / name

    def exists(self) -> bool:
        return self.path.is_dir()

    def state(self) -> PackageState:
        return PackageState(self, self.files())

    def files(self) -> list[Path]:
        if not self.exists():
            return []

        all_files: list[Path] = []
        for root, _, files in os.walk(self.path):
            relative_root = Path(root).relative_to(self.path)
            for file in files:
                all_files.append(relative_root / file)

        return all_files

    def __lt__(self, other: Any) -> bool:
        if not isinstance(other, Package):
            raise NotImplemented
        return self.name < other.name

    def __eq__(self, other: Any) -> bool:
        if not isinstance(other, Package):
            raise NotImplemented
        return self.name == other.name

    def __str__(self) -> str:
        return f'package "{self.name}"'

    def __repr__(self) -> str:
        return f"<Package name={self.name}>"


def get_all_packages() -> list[Package]:
    return sorted(Package(name) for name in os.listdir(PACKAGE_ROOT))


STOW_ARGS = ["-d", str(PACKAGE_ROOT), "-t", str(HOME)]


def stow_packages(pkgs: list[Package]):
    return subprocess.run(["stow"] + STOW_ARGS + ["-S"] + [pkg.name for pkg in pkgs])


def unstow_packages(pkgs: list[Package]):
    return subprocess.run(["stow"] + STOW_ARGS + ["-D"] + [pkg.name for pkg in pkgs])


def package_status(args: argparse.Namespace):
    pkg = Package(args.package)
    if not pkg.exists():
        print(f"Error: {pkg} does not exist")
        sys.exit(1)
    else:
        print(pkg.state())


def list_packages(args: argparse.Namespace):
    for pkg in get_all_packages():
        if args.verbose:
            print(pkg.state())
        else:
            print(pkg.name)


def install_packages(args: argparse.Namespace):
    ok = True

    if args.all:
        packages = get_all_packages()
    else:
        packages = [Package(name) for name in args.package]

    for pkg in packages:
        if not pkg.exists():
            print(f"Error: {pkg} does not exist")
            ok = False

    if not ok:
        sys.exit(1)

    notices: list[str] = []
    errors: list[str] = []
    queue: list[Package] = []

    for pkg in packages:
        state = pkg.state()

        if state.is_installed():
            notices.append(f"Skipped: {pkg} is already installed")
        elif state.is_empty():
            notices.append(f"Skipped: {pkg} is empty")
        elif not state.can_install():
            errors.append(f"Error: {pkg} cannot be installed\n{state}")
        else:
            queue.append(pkg)

    for msg in notices + errors:
        print(msg)

    if len(errors) > 0:
        sys.exit(1)

    if len(queue) == 0:
        return

    result = stow_packages(queue)
    if result.returncode != 0:
        print(f"Error: stow failed")
        print(f"Command: {' '.join(result.args)}")
        print(result.stderr)
    else:
        for pkg in queue:
            print(f"{pkg} installed")


def uninstall_packages(args: argparse.Namespace):
    ok = True

    if args.all:
        packages = get_all_packages()
    else:
        packages = [Package(name) for name in args.package]

    for pkg in packages:
        if not pkg.exists():
            print(f"Error: {pkg} does not exist")
            ok = False

    if not ok:
        sys.exit(1)

    notices: list[str] = []
    errors: list[str] = []
    queue: list[Package] = []

    for pkg in packages:
        state = pkg.state()

        if state.is_uninstalled():
            notices.append(f"Skipped: {pkg} is already uninstalled")
        elif state.is_empty():
            notices.append(f"Skipped: {pkg} is empty")
        else:
            queue.append(pkg)

    for msg in notices + errors:
        print(msg)

    if len(queue) == 0:
        return

    result = unstow_packages(queue)
    if result.returncode != 0:
        print(f"Error: stow failed")
        print(f"Command: {' '.join(result.args)}")
        print(result.stderr)
    else:
        for pkg in queue:
            print(f"{pkg} uninstalled")


def main():
    parser = argparse.ArgumentParser()
    subparsers = parser.add_subparsers(title="subcommands", metavar="", required=True)

    list_parser = subparsers.add_parser(
        "list", aliases=["ls"], help="List all packages"
    )
    list_parser.add_argument(
        "-v", "--verbose", help="Show package status", action="store_true"
    )
    list_parser.set_defaults(func=list_packages)

    status_parser = subparsers.add_parser("status", help="Show the status of a package")
    status_parser.add_argument("package", help="Name of the package")
    status_parser.set_defaults(func=package_status)

    install_parser = subparsers.add_parser(
        "install", help="Install packages into the home directory"
    )
    install_parser.add_argument("package", help="Package(s) to install", nargs="*")
    install_parser.add_argument(
        "--all", help="Install all packages", action="store_true"
    )
    install_parser.set_defaults(func=install_packages)

    uninstall_parser = subparsers.add_parser(
        "uninstall", help="Uninstall packages from the home directory"
    )
    uninstall_parser.add_argument("package", help="Package(s) to uninstall", nargs="*")
    uninstall_parser.add_argument(
        "--all", help="uninstall all packages", action="store_true"
    )
    uninstall_parser.set_defaults(func=uninstall_packages)

    args = parser.parse_args()
    args.func(args)


if __name__ == "__main__":
    main()
