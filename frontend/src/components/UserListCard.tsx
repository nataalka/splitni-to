import * as React from "react"
import { Card } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { UserListItem } from "@/components/UserListItem"
import type { User } from "@/types";

interface UserListCardProps {
  users: User[] | undefined;
  renderActions?: (user: User) => React.ReactNode;
  renderSubtext?: (user: User) => React.ReactNode;
}

export function UserListCard({
   users,
   renderActions,
   renderSubtext
 }: UserListCardProps) {
  return (
    <Card className="overflow-hidden p-2">
      <div className="flex flex-col gap-2">
        {users.map((user, index) => (
          <React.Fragment key={user.id}>
            <UserListItem
              user={user}
              actions={renderActions?.(user)}
              subtext={renderSubtext?.(user)}
            />
            {index < users.length - 1 && <Separator/>}
          </React.Fragment>
        ))}
      </div>
    </Card>
  )
}