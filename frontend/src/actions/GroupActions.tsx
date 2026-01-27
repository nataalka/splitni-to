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
          <Button variant="outline">
            <MoreHorizontal className="h-5 w-5" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-40">
          <DropdownMenuLabel className="text-xs font-bold text-zinc-400 uppercase tracking-widest">
            Group Options
          </DropdownMenuLabel>
          <DropdownMenuSeparator />

          <DropdownMenuItem
            onSelect={() => setShowEdit(true)}
          >
            <Settings2 className="h-4 w-4" />
            <span>Edit Group</span>
          </DropdownMenuItem>

          <DropdownMenuItem
            onSelect={() => setShowDelete(true)}
            variant="destructive"
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