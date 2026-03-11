import * as React from "react"
import { useForm, Controller } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import api from "@/lib/api"

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
import { UserPlus } from "lucide-react"

const friendSchema = z.object({
  email: z.string().email("Please enter a valid email address."),
})

type FriendValues = z.infer<typeof friendSchema>

export function AddFriendDialog() {
  const [open, setOpen] = React.useState(false)
  const queryClient = useQueryClient()

  const form = useForm<FriendValues>({
    resolver: zodResolver(friendSchema),
    defaultValues: { email: "" },
  })

  const mutation = useMutation({
    mutationFn: (data: FriendValues) => api.post("/friends", data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["friends"] })
      queryClient.invalidateQueries({ queryKey: ["friends", "pending"] })
      toast.success("Friend request sent!")
      setOpen(false)
      form.reset()
    },
    onError: (error: any) => {
      const message = error.response?.data?.error || "Failed to send request."
      toast.error(message)
    },
  })

  function onSubmit(data: FriendValues) {
    mutation.mutate(data)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="pinkPrimary">
          <UserPlus className="h-4 w-4" /> Add Friend
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[400px] rounded-3xl overflow-hidden gap-2">
        <DialogHeader>
          <DialogTitle className="text-2xl font-black">Add a Friend</DialogTitle>
          <DialogDescription>
            Enter your friend's email address to send them an invite.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <FieldGroup>
            <Controller
              name="email"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field data-invalid={fieldState.invalid}>
                  <FieldLabel>Email Address</FieldLabel>
                  <Input
                    {...field}
                    type="email"
                    placeholder="friend@example.com"
                    className="rounded-xl border-zinc-200 focus-visible:ring-pink-500"
                  />
                  {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                </Field>
              )}
            />
          </FieldGroup>
          <DialogFooter>
            <Button
              type="submit"
              variant="pinkPrimary"
              className="w-full sm:w-auto px-8 mt-4"
              disabled={mutation.isPending}
            >
              {mutation.isPending ? "Sending..." : "Send Invite"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}