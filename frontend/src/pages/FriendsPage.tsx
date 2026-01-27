import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import api from "@/lib/api"
import { Button } from "@/components/ui/button"
import { Check, X } from "lucide-react"
import * as React from "react"
import { toast } from "sonner"
import { UserListCard } from "@/components/UserListCard.tsx";
import { AddFriendDialog } from "@/dialogs/AddFriendDialog.tsx";


export default function FriendsPage() {
  const queryClient = useQueryClient()

  const {data: friends} = useQuery({
    queryKey: ["friends"], queryFn: () => api.get("/friends").then(res => res.data)
  })

  const {data: pending} = useQuery({
    queryKey: ["friends", "pending"], queryFn: () => api.get("/friends/pending").then(res => res.data)
  })

  const acceptMutation = useMutation({
    mutationFn: (requester_id: string) => api.put(`/friends/${requester_id}/accept`), onSuccess: () => {
      queryClient.invalidateQueries({queryKey: ["friends"]});
      toast.success("Friend added!");
    }
  });

  const rejectMutation = useMutation({
    mutationFn: (friend_id: string) => api.delete(`/friends/${friend_id}`), onSuccess: () => {
      queryClient.invalidateQueries({queryKey: ["friends", "pending"]});
      toast.error("Request removed");
    }
  });

  return (
    <div className="max-w-2xl mx-auto space-y-4">
      <div className="flex-col items-start justify-between gap-4">
        <div className="flex justify-between">
          <h1 className="text-3xl font-black text-zinc-900 tracking-tight">
            Friends
          </h1>
          <AddFriendDialog/>
        </div>
        <p className="text-zinc-500 text-sm">Add people to split expenses with them later.</p>
      </div>

      <div className="space-y-2">
        <h2 className="text-xs font-bold uppercase tracking-widest text-zinc-400 px-1">
          Your Friends
        </h2>
        <UserListCard
          users={friends}
          emptyMessage={"No friends yet."}
          renderActions={(user) => (
            <Button size="icon" variant="ghost" className="h-8 w-8 text-red-600 hover:bg-red-100"
                    onClick={() => rejectMutation.mutate(user.id)}>
              <X className="h-5 w-5"/>
            </Button>
          )}
        />
      </div>

      <div className="space-y-2">
        <h2 className="text-xs font-bold uppercase tracking-widest text-zinc-400 px-1">
          Pending Friend Requests
        </h2>
        <UserListCard
          users={pending}
          emptyMessage={"No pending requests"}
          renderActions={(user) => (
            <div>
              <Button size="icon" variant="ghost" className="h-8 w-8 text-green-600 hover:bg-green-100"
                      onClick={() => acceptMutation.mutate(user.id)}>
                <Check className="h-5 w-5"/>
              </Button>
              <Button size="icon" variant="ghost" className="h-8 w-8 text-red-600 hover:bg-red-100"
                      onClick={() => rejectMutation.mutate(user.id)}>
                <X className="h-5 w-5"/>
              </Button>
            </div>
          )}
        />
      </div>
    </div>)
}