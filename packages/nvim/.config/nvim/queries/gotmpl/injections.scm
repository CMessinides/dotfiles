;; extends

; Dynamic injection for gotmpl files (*.{ft}.tmpl)
; See init.lua for `inject-go-tmpl!` definition.
; Also see https://github.com/nvim-treesitter/nvim-treesitter/discussions/1917#discussioncomment-10714144
((text) @injection.content
  (#inject-go-tmpl!)
  (#set! injection.combined))
