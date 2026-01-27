import { useQuery } from "@tanstack/react-query"
import api from "@/lib/api"
import type { Group } from "@/types"
import { GroupListCard } from "@/components/GroupListCard.tsx";
import { GroupFormDialog } from "@/dialogs/GroupFormDialog.tsx";

export default function GroupsPage() {
  const {data: groups, isLoading, error} = useQuery<Group[]>({
    queryKey: ["groups"], queryFn: async () => {
      const response = await api.get("/groups")
      return response.data
    },
  })

  if (isLoading) return <div className="p-8 text-center">Loading groups...</div>
  if (error) return <div className="p-8 text-center text-destructive">Failed to load groups.</div>

  return (
    <div className="max-w-2xl mx-auto space-y-4">
      <div className="flex-col items-start justify-between gap-4">
        <div className="flex justify-between">
          <h1 className="text-3xl font-black text-zinc-900 tracking-tight">Groups</h1>
          <GroupFormDialog/>
        </div>
        <p className="text-zinc-500 text-sm mt-1">Manage your shared expenses and groups.</p>
      </div>

      <section className="space-y-2">
        <h2 className="text-xs font-bold uppercase tracking-widest text-zinc-400 px-1">
          Your active circles
        </h2>
        <GroupListCard groups={groups}/>
      </section>
    </div>)
}