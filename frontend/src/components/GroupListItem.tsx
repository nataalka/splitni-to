import type { Group } from "@/types";
import { ChevronRight, Users } from "lucide-react";
import { Link } from "react-router-dom";
import * as React from "react";

interface GroupListItemProps {
  group: Group;
}

export function GroupListItem({ group }: GroupListItemProps) {
  return (
    <Link
      to={`/groups/${group.id}`}
      className="group flex items-center justify-between p-4 hover:bg-zinc-50 transition-colors"
    >
      <div className="flex items-center gap-4">
        <div
          className="h-12 w-12 rounded-2xl bg-pink-50 flex items-center justify-center text-pink-600 border border-pink-100 group-hover:bg-pink-100 transition-colors">
          <Users className="h-6 w-6"/>
        </div>

        <div className="min-w-0">
          <h3 className="font-bold text-zinc-900 leading-tight">
            {group.name}
          </h3>
          <p className="text-xs text-zinc-500 truncate mt-1 max-w-[200px] sm:max-w-xs">
            {group.description || "No description provided."}
          </p>
        </div>
      </div>

      <div className="flex items-center gap-2">
        <ChevronRight className="h-5 w-5 text-zinc-300 group-hover:text-pink-500 transition-colors"/>
      </div>
    </Link>
  )
}