import { useQuery } from "@tanstack/react-query"
import api from "@/lib/api"
import type { Group } from "@/types"
import { Card, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Plus, Users } from "lucide-react"

export default function GroupsPage() {
  const { data: groups, isLoading, error } = useQuery<Group[]>({
    queryKey: ["groups"],
    queryFn: async () => {
      const response = await api.get("/groups")
      return response.data
    },
  })

  if (isLoading) return <div className="p-8 text-center">Loading groups...</div>
  if (error) return <div className="p-8 text-center text-destructive">Failed to load groups.</div>

  return (
    <div className="max-w-5xl mx-auto p-8 space-y-8">
      <div className="flex items-top justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">My Groups</h1>
          <p className="text-muted-foreground">Manage your shared expenses and groups.</p>
        </div>
        <Button className="bg-pink-600 hover:bg-pink-700">
          <Plus className="h-4 w-4" /> New Group
        </Button>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {groups?.map((group) => (
          <Card
            key={group.id}
            className="hover:border-pink-500 transition-all cursor-pointer shadow-sm"
          >
            <CardHeader className="flex items-center gap-4">
              <div className="p-2 bg-pink-100 rounded-lg">
                <Users className="h-6 w-6 text-pink-600" />
              </div>
              <div>
                <CardTitle>{group.name}</CardTitle>
                <CardDescription className="line-clamp-2">
                  {group.description || "No description provided."}
                </CardDescription>
              </div>
            </CardHeader>
          </Card>
        ))}

        {groups?.length === 0 && (
          <div className="col-span-full border-2 border-dashed rounded-xl p-12 text-center space-y-3">
            <p className="text-muted-foreground">You don't have any groups yet.</p>
            <Button variant="outline">Create your first group</Button>
          </div>
        )}
      </div>
    </div>
  )
}