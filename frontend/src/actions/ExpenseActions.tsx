import * as React from "react"
import { MoreHorizontal, Settings2, Trash2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem, DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { DeleteExpenseDialog } from "@/dialogs/DeleteExpenseDialog"
import { ExpenseFormDialog } from "@/dialogs/ExpenseFormDialog"
import type { ExpenseDetailed, User } from "@/types";

interface ExpenseActionsProps {
  groupId: string
  expense: ExpenseDetailed
  members: User[]
}

export function ExpenseActions({ groupId, expense, members }: ExpenseActionsProps) {
  const [showDelete, setShowDelete] = React.useState(false)
  const [showEdit, setShowEdit] = React.useState(false)

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline">
            <MoreHorizontal className="h-5 w-5" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-40">
          <DropdownMenuLabel className="text-xs font-bold text-zinc-400 uppercase tracking-widest">
            Expense Options
          </DropdownMenuLabel>
          <DropdownMenuSeparator />

          <DropdownMenuItem
            onSelect={(e) => {
              e.preventDefault();
              setShowEdit(true);
            }}
          >
            <Settings2 className="h-4 w-4" />
            <span>Edit Expense</span>
          </DropdownMenuItem>

          <DropdownMenuItem
            onSelect={(e) => {
              e.preventDefault();
              setShowDelete(true);
            }}
            variant="destructive"
          >
            <Trash2 className="h-4 w-4" />
            <span>Delete</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {showEdit && (
        <ExpenseFormDialog
          key={showEdit ? `edit-${expense.id}` : 'add'}
          groupId={groupId}
          members={members}
          expense={expense}
          open={showEdit}
          onOpenChange={setShowEdit}
        />
      )}

      {showDelete && (
        <DeleteExpenseDialog
          groupId={groupId}
          expenseId={expense.id}
          open={showDelete}
          onOpenChange={setShowDelete}
        />
      )}
    </>
  )
}