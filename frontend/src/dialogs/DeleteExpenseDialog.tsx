import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { useNavigate } from "react-router-dom"
import { toast } from "sonner"
import api from "@/lib/api"

interface DeleteExpenseDialogProps {
  groupId: string
  expenseId: string
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function DeleteExpenseDialog({ groupId, expenseId, open, onOpenChange }: DeleteExpenseDialogProps) {
  const queryClient = useQueryClient()
  const navigate = useNavigate()

  const mutation = useMutation({
    mutationFn: () => api.delete(`groups/${groupId}/expenses/${expenseId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["group", groupId, "expenses"] })
      toast.success("Expense deleted.")
      navigate(`/groups/${groupId}`)
    },
    onError: () => toast.error("Failed to delete expense."),
  })

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent className="rounded-[32px]">
        <AlertDialogHeader>
          <AlertDialogTitle className="text-xl font-black">Delete expense?</AlertDialogTitle>
          <AlertDialogDescription>
            This will remove this record and update everyone's balances. This action is permanent.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel className="rounded-xl">Cancel</AlertDialogCancel>
          <AlertDialogAction
            onClick={(e) => {
              e.preventDefault()
              mutation.mutate()
            }}
            className="bg-red-600 hover:bg-red-700 rounded-xl"
          >
            Delete
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}