(setq +org-capture-todo-file "inbox.org")

(after! (org-capture)
  (add-to-list 'org-capture-templates
               '("m" "Meeting"
                 entry (file+datetree "agenda.org")
                 "* %U %? %^G\n [[%L][Meeting Notes]]"
                 :clock-in t
                 :clock-resume t
                 :empty-lines 1))
  (add-to-list 'org-capture-templates
               '("s" "Story"
                 entry (file+headline "work.org" "Uncategorized")
                 "* TODO %?\n\n [[https://app.clickup.com/t/%^{Task ID}][Task %\\1]]"
                 :clock-in t
                 :clock-keep t
                 :jump-to-captured t
                 :empty-lines 1))
  (add-to-list 'org-capture-templates
               '("c" "Call"
                 entry (file+headline +org-capture-todo-file "Calls")
                 "* %U Call with %?"
                 :empty-lines 1)))
