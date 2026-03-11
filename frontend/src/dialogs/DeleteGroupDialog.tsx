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

interface DeleteGroupDialogProps {
  groupId: string
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function DeleteGroupDialog({ groupId, open, onOpenChange }: DeleteGroupDialogProps) {
  const queryClient = useQueryClient()
  const navigate = useNavigate()

  const mutation = useMutation({
    mutationFn: () => api.delete(`/groups/${groupId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["groups"] })
      toast.success("Group deleted successfully.")
      navigate("/groups")
    },
    onError: () => {
      toast.error("Could not delete group. Please try again.")
      onOpenChange(false)
    },
  })

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent className="rounded-[32px]">
        <AlertDialogHeader>
          <AlertDialogTitle className="text-xl font-black">Delete this group?</AlertDialogTitle>
          <AlertDialogDescription className="text-zinc-500">
            This action cannot be undone. All expenses and balances for this group will be permanently removed.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter className="gap-2">
          <AlertDialogCancel className="rounded-xl border-zinc-100">Cancel</AlertDialogCancel>
          <AlertDialogAction
            onClick={(e) => {
              e.preventDefault()
              mutation.mutate()
            }}
            className="bg-red-600 hover:bg-red-700 text-white rounded-xl"
            disabled={mutation.isPending}
          >
            {mutation.isPending ? "Deleting..." : "Yes, Delete Group"}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}