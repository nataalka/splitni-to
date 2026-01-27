import * as React from "react"
import { Card } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { UserListItem } from "@/components/UserListItem"
import type { User } from "@/types";

interface UserListCardProps {
  users: User[] | undefined;
  emptyMessage?: string;
  renderActions: React.ReactNode;
}

export function UserListCard({ users, emptyMessage = "No people found.", renderActions }: UserListCardProps) {
  return (
    <Card className="overflow-hidden p-2">
      {users && users.length > 0 ? (
        <div className="flex flex-col gap-2">
          {users.map((user, index) => (
            <React.Fragment key={user.id}>
              <UserListItem
                user={user}
                actions={renderActions}
              />
              {index < users.length - 1 && <Separator/>}
            </React.Fragment>
          ))}
        </div>
      ) : (
        <div className="p-8 text-center">
          <p className="text-sm text-zinc-500 italic">{emptyMessage}</p>
        </div>
      )}
    </Card>
  )
}