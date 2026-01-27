import * as React from "react"
import type { User } from "@/types";

interface UserListItemProps {
  user: User;
  actions?: React.ReactNode;
  subtext?: React.ReactNode;
  className?: string;
}

export function UserListItem({ user, actions, subtext, className = "" }: UserListItemProps) {
  const initials = `${user.name?.[0] || ""}${user.surname?.[0] || ""}`.toUpperCase();

  return (
    <div className={`flex items-center justify-between p-2 rounded-xl hover:bg-zinc-50 transition-colors group ${className}`}>
      <div className="flex items-center gap-3">
        <div className="h-9 w-9 rounded-full bg-pink-100 text-pink-700 flex items-center justify-center font-bold text-xs shrink-0">
          {initials}
        </div>

        <div className="min-w-0">
          <p className="text-sm font-semibold leading-none text-zinc-900 truncate">
            {user.name} {user.surname}
          </p>
          <p className="text-xs text-muted-foreground truncate mt-1">
            {subtext ? subtext : <span>{user.email}</span>}
          </p>
        </div>
      </div>

      {actions && (
        <div className="flex items-center gap-1">
          {actions}
        </div>
      )}
    </div>
  )
}