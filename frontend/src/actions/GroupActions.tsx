import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Button } from "@/components/ui/button"
import { MoreHorizontal, Settings2, Trash2 } from "lucide-react"
import { GroupFormDialog } from "@/dialogs/GroupFormDialog"
import { DeleteGroupDialog } from "@/dialogs/DeleteGroupDialog"
import type { Group } from "@/types"
import * as React from "react"

export function GroupActions({ group }: { group: Group }) {
  const [showEdit, setShowEdit] = React.useState(false)
  const [showDelete, setShowDelete] = React.useState(false)

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon" className="h-9 w-9 text-zinc-400 rounded-full">
            <MoreHorizontal className="h-5 w-5" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-48 rounded-2xl p-2">
          <DropdownMenuLabel className="text-xs font-bold text-zinc-400 uppercase tracking-widest px-2 py-1.5">
            Group Options
          </DropdownMenuLabel>
          <DropdownMenuSeparator />

          <DropdownMenuItem
            onSelect={() => setShowEdit(true)}
            className="rounded-xl gap-2 cursor-pointer py-2.5"
          >
            <Settings2 className="h-4 w-4 text-zinc-500" />
            <span>Edit Group</span>
          </DropdownMenuItem>

          <DropdownMenuItem
            onSelect={() => setShowDelete(true)}
            className="rounded-xl gap-2 cursor-pointer py-2.5 text-red-600 focus:text-red-600 focus:bg-red-50"
          >
            <Trash2 className="h-4 w-4" />
            <span>Delete Group</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {showEdit && (
        <GroupFormDialog
          group={group}
          open={showEdit}
          onOpenChange={setShowEdit}
        />
      )}
      {showDelete && (
        <DeleteGroupDialog
          groupId={group.id}
          open={showDelete}
          onOpenChange={setShowDelete}
        />
      )}
    </>
  )
}