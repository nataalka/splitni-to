import * as React from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import api from "@/lib/api"
import { toast } from "sonner"
import { UserPlus } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent, DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { ScrollArea } from "@/components/ui/scroll-area.tsx";
import { EmptyBox } from "@/components/EmptyBox.tsx";

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
        <Button variant="pinkOutline">
          <UserPlus className="h-4 w-4"/> Add Member
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[450px] rounded-3xl overflow-hidden gap-2">
        <DialogHeader className="flex justify-start gap-2">
          <DialogTitle className="text-2xl font-black">Add Friend to Group</DialogTitle>
          <DialogDescription>Your Friends</DialogDescription>
        </DialogHeader>

        <ScrollArea>
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
            <EmptyBox description="You don't have any friends yet."/>
          )}
        </ScrollArea>
      </DialogContent>
    </Dialog>
  )
}