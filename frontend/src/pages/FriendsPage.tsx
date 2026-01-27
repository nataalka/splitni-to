import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import api from "@/lib/api"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardHeader, CardContent } from "@/components/ui/card"
import { Users, User, Check, X } from "lucide-react"
import { useState } from "react"
import { toast } from "sonner"

export default function FriendsPage() {
  const [email, setEmail] = useState("")
  const queryClient = useQueryClient()

  const {data: friends} = useQuery({
    queryKey: ["friends"],
    queryFn: () => api.get("/friends").then(res => res.data)
  })

  const {data: pending} = useQuery({
    queryKey: ["friends", "pending"],
    queryFn: () => api.get("/friends/pending").then(res => res.data)
  })

  const addFriendMutation = useMutation({
    mutationFn: (email: string) => api.post("/friends", {email}),
    onSuccess: () => {
      toast.success("Request sent!")
      setEmail("")
      queryClient.invalidateQueries({queryKey: ["friends"]})
    }
  })

  const acceptMutation = useMutation({
    mutationFn: (requester_id: string) => api.put(`/friends/${requester_id}/accept`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["friends"] });
      toast.success("Friend added!");
    }
  });

  const rejectMutation = useMutation({
    mutationFn: (friend_id: string) => api.delete(`/friends/${friend_id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["friends", "pending"] });
      toast.error("Request removed");
    }
  });

  return (
    <div className="max-w-2xl mx-auto p-6 space-y-4">
      <Card>
        <CardHeader>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <Users className="text-pink-600"/> Friends
          </h1>
          <p className="text-muted-foreground">Add people to split expenses with them later.</p>
        </CardHeader>
      </Card>

      <CardContent className="flex gap-4 w-full">
        <Input
          placeholder="Enter friend's email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <Button onClick={() => addFriendMutation.mutate(email)}>
          Add Friend
        </Button>
      </CardContent>

      <div className="space-y-4">
        <h2 className="font-semibold text-lg">Your Friends</h2>
        {friends && friends.length > 0 ? (
          friends.map((friend: any) => (
            <Card key={friend.id}>
              <CardContent className="gap-4 flex justify-between items-center">
                <div className="gap-4 flex justify-start items-center">
                  <div className="p-2 bg-pink-100 rounded-lg">
                    <User className="h-6 w-6 text-pink-600"/>
                  </div>
                  <div className="flex-row items-center gap-2">
                    <p className="font-medium">{friend.name} {friend.surname}</p>
                    <p className="text-sm text-gray-500">{friend.email}</p>
                  </div>
                </div>
                <Button size="icon" variant="ghost" className="h-8 w-8 text-red-600 hover:bg-red-100"
                        onClick={() => rejectMutation.mutate(friend.id)}>
                  <X className="h-5 w-5" />
                </Button>
              </CardContent>
            </Card>
          ))
        ) : (
          <p className="text-sm text-gray-500">No friends yet.</p>
        )}
      </div>

      <div className="space-y-4">
        <h2 className="font-semibold text-lg">Pending Friend Requests</h2>
        {pending && pending.length > 0 ? (
          pending.map((p: any) => (
            <Card key={p.id}>
              <CardContent className="gap-4 flex justify-between items-center">
                <div className="gap-4 flex justify-start items-center">
                  <div className="p-2 bg-pink-100 rounded-lg">
                    <User className="h-6 w-6 text-pink-600"/>
                  </div>
                  <div className="flex-row items-center gap-2">
                    <p className="font-medium">{p.name} {p.surname}</p>
                    <p className="text-sm text-gray-500">{p.email}</p>
                  </div>
                </div>
                <div className="flex gap-2">
                  <Button size="icon" variant="ghost" className="h-8 w-8 text-green-600 hover:bg-green-100"
                          onClick={() => acceptMutation.mutate(p.id)}>
                    <Check className="h-5 w-5" />
                  </Button>
                  <Button size="icon" variant="ghost" className="h-8 w-8 text-red-600 hover:bg-red-100"
                          onClick={() => rejectMutation.mutate(p.id)}>
                    <X className="h-5 w-5" />
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))
        ) : (
          <p className="text-sm text-gray-500">No pending requests.</p>
        )}
      </div>
    </div>
  )
}