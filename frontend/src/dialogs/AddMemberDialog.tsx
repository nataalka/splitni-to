import * as React from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import api from "@/lib/api"
import { toast } from "sonner"
import { UserPlus } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { ScrollArea } from "@/components/ui/scroll-area.tsx";

export function AddMemberDialog({groupId}: { groupId: string }) {
  const [open, setOpen] = React.useState(false)
  const queryClient = useQueryClient()

  const {data: friends, isLoading} = useQuery({
    queryKey: ["friends"],
    queryFn: () => api.get("/friends").then(res => res.data)
  })

  const mutation = useMutation({
    mutationFn: async (userId: string) => {
      return api.post(`/groups/${groupId}/members`, {user_id: userId})
    },
    onSuccess: () => {
      queryClient.invalidateQueries({queryKey: ["group", groupId]})
      toast.success("Member added to group!")
      setOpen(false)
      setEmail("")
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || "User not found or already in group.")
    }
  })

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="pinkOutline" size="sm">
          <UserPlus className="h-4 w-4 mr-2 text-pink-600"/> Add Friend
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[400px]">
        <DialogHeader>
          <DialogTitle>Add Friend to Group</DialogTitle>
        </DialogHeader>

        <div>
          <h3 className="text-sm font-medium mb-3 text-muted-foreground">Your Friends</h3>

          <ScrollArea className="h-[300px]">
            {isLoading ? (
              <p className="text-sm text-center py-4">Loading friends...</p>
            ) : friends && friends.length > 0 ? (
              <div className="space-y-2">
                {friends.map((friend: any) => (
                  <div
                    key={friend.id}
                    className="flex items-center justify-between p-3 border rounded-xl hover:border-pink-300 hover:bg-pink-50/50 transition-all group"
                  >
                    <div className="flex items-center gap-3">
                      <div
                        className="h-10 w-10 rounded-full bg-pink-100 flex items-center justify-center text-pink-700 font-bold text-sm">
                        {friend.name[0]}{friend.surname[0]}
                      </div>
                      <div>
                        <p className="font-medium text-sm leading-none">{friend.name} {friend.surname}</p>
                        <p className="text-xs text-muted-foreground">{friend.email}</p>
                      </div>
                    </div>
                    <Button
                      size="sm"
                      variant="ghost"
                      className="text-pink-600 hover:text-pink-700 hover:bg-pink-100"
                      disabled={mutation.isPending}
                      onClick={() => mutation.mutate(friend.id)}
                    >
                      {mutation.isPending ? "Adding..." : "Add"}
                    </Button>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-10">
                <p className="text-sm text-muted-foreground">You don't have any friends yet.</p>
              </div>
            )}
          </ScrollArea>
        </div>
      </DialogContent>
    </Dialog>
  )
}