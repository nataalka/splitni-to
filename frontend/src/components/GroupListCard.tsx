import * as React from "react"
import { Card } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import type { Group } from "@/types"
import { GroupListItem } from "@/components/GroupListItem.tsx";

interface GroupListCardProps {
  groups: Group[] | undefined;
  onGroupClick?: (id: string) => void;
  as?: "link" | "button";
}

export function GroupListCard({ groups, onGroupClick, as }: GroupListCardProps) {
  return (
    <Card className="overflow-hidden p-2">
      <div className="flex flex-col gap-2">
        {groups?.map((group, index) => (
          <React.Fragment key={group.id}>
            <GroupListItem
              group={group}
              onClick={onGroupClick}
              as={as}/>
            {index < groups.length - 1 && <Separator/>}
          </React.Fragment>
        ))}
      </div>
    </Card>
  )
}