import * as React from "react"
import { useForm, Controller } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import api from "@/lib/api"
import { Plus, Receipt } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field"
import type { Group } from "@/types"

const groupSchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters."),
  description: z.string().max(100, "Description is too long.").optional(),
})

type GroupValues = z.infer<typeof groupSchema>

interface GroupFormDialogProps {
  group?: Group; // If provided, we are in edit mode
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}

export function GroupFormDialog({ group, open: externalOpen, onOpenChange: setExternalOpen }: GroupFormDialogProps) {
  const [internalOpen, setInternalOpen] = React.useState(false)
  const queryClient = useQueryClient()

  const isEdit = !!group
  const isOpen = externalOpen ?? internalOpen
  const setIsOpen = setExternalOpen ?? setInternalOpen

  const form = useForm<GroupValues>({
    resolver: zodResolver(groupSchema),
    defaultValues: {
      name: group?.name || "",
      description: group?.description || "",
    },
  })

  React.useEffect(() => {
    if (group) {
      form.reset({ name: group.name, description: group.description || "" })
    }
  }, [group, form])

  const mutation = useMutation({
    mutationFn: (data: GroupValues) =>
      isEdit ? api.put(`/groups/${group.id}`, data) : api.post("/groups", data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["groups"] })
      if (isEdit) queryClient.invalidateQueries({ queryKey: ["group", group.id] })

      toast.success(isEdit ? "Group updated!" : "Group created!")
      setIsOpen(false)
      if (!isEdit) form.reset()
    },
    onError: () => toast.error("An error occurred. Please try again."),
  })

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      {!externalOpen && (
        <DialogTrigger asChild>
          <Button variant="pinkPrimary">
            <Plus className="h-4 w-4" /> New Group
          </Button>
        </DialogTrigger>
      )}

      <DialogContent className="sm:max-w-[425px] rounded-3xl]">
        <DialogHeader className="flex-row justify-start gap-2">
          <div
            className="h-12 w-12 rounded-2xl bg-pink-50 flex items-center justify-center text-pink-600 mb-2 border border-pink-100">
            <Receipt className="h-6 w-6"/>
          </div>
          <div className="flex-col">
            <DialogTitle className="text-2xl font-black">
              {isEdit ? "Edit Group" : "Create Group"}
            </DialogTitle>
            <DialogDescription>
              {isEdit ? "Update group name or description." : "Start tracking expenses with friends."}
            </DialogDescription>
          </div>
        </DialogHeader>

        <form onSubmit={form.handleSubmit((data) => mutation.mutate(data))}>
          <FieldGroup className="py-4">
            <Controller
              name="name"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field data-invalid={fieldState.invalid}>
                  <FieldLabel>Group Name</FieldLabel>
                  <Input {...field} placeholder="e.g. Ski Trip 2026" className="rounded-xl" />
                  {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                </Field>
              )}
            />
            <Controller
              name="description"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field data-invalid={fieldState.invalid}>
                  <FieldLabel>Description</FieldLabel>
                  <Input {...field} placeholder="Optional details..." className="rounded-xl" />
                  {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                </Field>
              )}
            />
          </FieldGroup>

          <DialogFooter>
            <Button
              type="submit"
              variant="pinkPrimary"
              className="w-full rounded-xl"
              disabled={mutation.isPending}
            >
              {mutation.isPending ? "Saving..." : isEdit ? "Save Changes" : "Create Group"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}