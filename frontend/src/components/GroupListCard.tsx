import * as React from "react"
import { Card } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { Users } from "lucide-react"
import type { Group } from "@/types"
import { GroupListItem } from "@/components/GroupListItem.tsx";

interface GroupListCardProps {
  groups: Group[] | undefined;
}

export function GroupListCard({ groups }: GroupListCardProps) {
  return (
    <Card className="overflow-hidden p-2">
      {groups && groups.length > 0 ? (
        <div className="flex flex-col gap-2">
          {groups.map((group, index) => (
            <React.Fragment key={group.id}>
              <GroupListItem group={group}/>
              {index < groups.length - 1 && <Separator/>}
            </React.Fragment>
          ))}
        </div>
      ) : (
        <div className="py-16 text-center">
          <div className="bg-zinc-50 h-16 w-16 rounded-full flex items-center justify-center mx-auto mb-4">
            <Users className="h-8 w-8 text-zinc-200" />
          </div>
          <p className="text-zinc-400 text-sm italic px-6">
            You don't have any groups yet. Start by creating one!
          </p>
        </div>
      )}
    </Card>
  )
}