import * as React from "react"
import { Controller, useForm } from "react-hook-form"
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
import { Checkbox } from "@/components/ui/checkbox"
import type { CreateExpenseRequest, User } from "@/types"
import { UserListCard } from "@/components/UserListCard.tsx";

const expenseSchema = z.object({
  description: z.string().min(3, "Description is too short."),
  amount: z.string().refine((val) => !isNaN(Number(val)) && Number(val) > 0, "Invalid amount"),
  payer_id: z.string().uuid(),
  selectedMembers: z.array(z.string()).min(1, "Select at least one person to split with."),
})

type ExpenseValues = z.infer<typeof expenseSchema>

interface AddExpenseDialogProps {
  groupId: string;
  members: User[];
}

export function AddExpenseDialog({groupId, members}: AddExpenseDialogProps) {
  const [open, setOpen] = React.useState(false)
  const queryClient = useQueryClient()

  const form = useForm<ExpenseValues>({
    resolver: zodResolver(expenseSchema), defaultValues: {
      description: "", amount: "", payer_id: members[0]?.id || "", selectedMembers: members.map(m => m.id),
    },
  })

  const watchedAmount = form.watch("amount")
  const watchedMembers = form.watch("selectedMembers")

  const mutation = useMutation({
    mutationFn: (data: CreateExpenseRequest) => api.post(`/groups/${groupId}/expenses`, data), onSuccess: () => {
      queryClient.invalidateQueries({queryKey: ["group", groupId, "expenses"]})
      queryClient.invalidateQueries({queryKey: ["group", groupId, "balances"]})
      toast.success("Expense added!")
      setOpen(false)
      form.reset()
    },
  })

  const calculatedSplits = React.useMemo(() => {
    const totalAmount = parseFloat(watchedAmount);
    const count = watchedMembers.length

    if (isNaN(totalAmount) || count === 0) return []

    const totalInCents = Math.round(totalAmount * 100);
    const baseShareInCents = Math.floor(totalInCents / count);
    let remainderInCents = totalInCents % count;

    return watchedMembers.map((userId) => {
      let finalShareInCents = baseShareInCents;

      if (remainderInCents > 0) {
        finalShareInCents += 1;
        remainderInCents -= 1;
      }

      return {
        user_id: userId,
        amount: (finalShareInCents / 100).toFixed(2)
      };
    });
  }, [watchedAmount, watchedMembers])

  function onSubmit(values: ExpenseValues) {
    const payload: CreateExpenseRequest = {
      description: values.description,
      amount: values.amount,
      currency: "EUR",
      payer_id: values.payer_id,
      splits: calculatedSplits
    }
    mutation.mutate(payload)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="pinkPrimary">
          <Plus className="h-4 w-4"/> Add Expense
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[450px] rounded-3xl overflow-hidden">
        <DialogHeader className="flex-row justify-start gap-2">
          <div
            className="h-12 w-12 rounded-2xl bg-pink-50 flex items-center justify-center text-pink-600 mb-2 border border-pink-100">
            <Receipt className="h-6 w-6"/>
          </div>
          <div className="flex-col">
            <DialogTitle className="text-2xl font-black">Add Expense</DialogTitle>
            <DialogDescription>Split a new bill with your group.</DialogDescription>
          </div>
        </DialogHeader>

        <form onSubmit={form.handleSubmit(onSubmit)}>
          <FieldGroup>
            <div className="grid grid-cols-3 gap-4">
              <div className="col-span-2">
                <Controller
                  name="description"
                  control={form.control}
                  render={({field, fieldState}) => (<Field data-invalid={fieldState.invalid}>
                    <FieldLabel>Description</FieldLabel>
                    <Input {...field} placeholder="Grocery shopping" className="rounded-xl"/>
                    {fieldState.invalid && <FieldError errors={[fieldState.error]}/>}
                  </Field>)}
                />
              </div>
              <div className="col-span-1">
                <Controller
                  name="amount"
                  control={form.control}
                  render={({field, fieldState}) => (<Field data-invalid={fieldState.invalid}>
                    <FieldLabel>Amount</FieldLabel>
                    <div className="relative">
                      <Input {...field} placeholder="0.00" className="rounded-xl pr-7"/>
                      <span className="absolute right-3 top-2.5 text-zinc-400 text-sm">€</span>
                    </div>
                    {fieldState.invalid && <FieldError errors={[fieldState.error]}/>}
                  </Field>)}
                />
              </div>
            </div>

            <div className="pt-2">
              <FieldLabel className="flex justify-between items-center">
                Split with
                <span className="text-[10px] font-bold text-pink-600 bg-pink-50 px-2 py-0.5 rounded-full uppercase">
                  Split Equally
                </span>
              </FieldLabel>
              <div className="max-h-[250px] overflow-y-auto">
                <UserListCard
                  users={members}
                  renderSubtext={(user) => {
                    const splitData = calculatedSplits.find(s => s.user_id === user.id);
                    const share = splitData?.amount

                    return share ? (
                      <span className="text-[11px] text-pink-600 font-medium tracking-wide">
                          owes {share} €
                      </span>
                    ) : null;
                  }}
                  renderActions={(user) => (
                    <Controller
                      name="selectedMembers"
                      control={form.control}
                      render={({field}) => (
                        <Checkbox
                          checked={field.value.includes(user.id)}
                          onCheckedChange={(checked) => {
                            const newValue = checked
                              ? [...field.value, user.id]
                              : field.value.filter(id => id !== user.id);
                            field.onChange(newValue);
                          }}
                        />
                      )}
                    />
                  )}
                />
              </div>

              {form.formState.errors.selectedMembers && (
                <p className="text-xs text-red-500 font-medium">{form.formState.errors.selectedMembers.message}</p>
              )}
            </div>
          </FieldGroup>

          <DialogFooter>
            <Button type="submit" variant="pinkPrimary" disabled={mutation.isPending}>
              {mutation.isPending ? "Adding..." : "Save Expense"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>)
}