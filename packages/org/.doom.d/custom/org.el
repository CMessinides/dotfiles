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
               '("c" "Call"
                 entry (file+headline +org-capture-todo-file "Calls")
                 "* %U Call with %?"
                 :empty-lines 1)))
