import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import api from "@/lib/api"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Check, UserPlus, Users, X } from "lucide-react"
import * as React from "react"
import { useState } from "react"
import { toast } from "sonner"
import { UserListCard } from "@/components/UserListCard.tsx";


export default function FriendsPage() {
  const [email, setEmail] = useState("")
  const queryClient = useQueryClient()

  const {data: friends} = useQuery({
    queryKey: ["friends"], queryFn: () => api.get("/friends").then(res => res.data)
  })

  const {data: pending} = useQuery({
    queryKey: ["friends", "pending"], queryFn: () => api.get("/friends/pending").then(res => res.data)
  })

  const addFriendMutation = useMutation({
    mutationFn: (email: string) => api.post("/friends", {email}), onSuccess: () => {
      toast.success("Request sent!")
      setEmail("")
      queryClient.invalidateQueries({queryKey: ["friends"]})
    }
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

  return (<div className="max-w-2xl mx-auto p-4 space-y-4">
      <section className="space-y-4">
        <div>
          <h1 className="text-2xl font-black text-zinc-900 flex items-center gap-2">
            <Users className="text-pink-600 h-6 w-6"/> Friends
          </h1>
          <p className="text-zinc-500 text-sm">Add people to split expenses with them later.</p>
        </div>

        <div
          className="flex items-center w-full rounded-2xl border border-zinc-200 bg-white focus-within:ring-2 focus-within:ring-pink-500 focus-within:border-transparent transition-all">
          <Input
            placeholder="friend@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="border-0 focus-visible:ring-0 focus-visible:ring-offset-0 bg-transparent h-10 px-4"
          />
          <Button
            onClick={() => addFriendMutation.mutate(email)}
            disabled={!email || addFriendMutation.isPending}
            className="h-10 rounded-l-none rounded-r-2xl bg-pink-600 hover:bg-pink-700 px-4 shadow-none"
          >
            <UserPlus className="h-4 w-4"/>
            Add
          </Button>
        </div>
      </section>

      <div className="space-y-4">
        <h2 className="font-semibold text-lg">Your Friends</h2>
        <UserListCard
          users={friends}
          emptyMessage={"No friends yet."}
          renderActions={
            <Button size="icon" variant="ghost" className="h-8 w-8 text-red-600 hover:bg-red-100"
                           onClick={() => rejectMutation.mutate(friend.id)}>
            <X className="h-5 w-5"/>
          </Button>}
        />
      </div>

      <div className="space-y-4">
        <h2 className="font-semibold text-lg">Pending Friend Requests</h2>
        <UserListCard
          users={pending}
          emptyMessage={"No pending requests"}
          renderActions={
          <div>
            <Button size="icon" variant="ghost" className="h-8 w-8 text-green-600 hover:bg-green-100"
                    onClick={() => acceptMutation.mutate(p.id)}>
              <Check className="h-5 w-5"/>
            </Button>
            <Button size="icon" variant="ghost" className="h-8 w-8 text-red-600 hover:bg-red-100"
                    onClick={() => rejectMutation.mutate(p.id)}>
              <X className="h-5 w-5"/>
            </Button>
          </div>
        }
        />
      </div>
    </div>)
}