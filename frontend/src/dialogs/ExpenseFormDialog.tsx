import * as React from "react"
import { Controller, useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
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
import type { CreateExpenseRequest, ExpenseDetailed, User } from "@/types"
import { UserListCard } from "@/components/UserListCard.tsx";
import { Switch } from "@/components/ui/switch.tsx";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select.tsx";

const expenseSchema = z.object({
  description: z.string().min(3, "Description is too short."),
  amount: z.string().refine((val) => !isNaN(Number(val)) && Number(val) > 0, "Invalid amount"),
  payer_id: z.string().uuid(),
  selectedMembers: z.array(z.string()).min(1, "Select at least one person to split with."),
})

type ExpenseValues = z.infer<typeof expenseSchema>

interface ExpenseFormDialogProps {
  groupId: string;
  expense?: ExpenseDetailed;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}

export function ExpenseFormDialog({
   groupId,
   expense,
   open: externalOpen,
   onOpenChange: setExternalOpen
 }: ExpenseFormDialogProps) {
  const [internalOpen, setInternalOpen] = React.useState(false)
  const isEdit = !!expense;
  const isOpen = externalOpen !== undefined ? externalOpen : internalOpen
  const setIsOpen = setExternalOpen !== undefined ? setExternalOpen : setInternalOpen

  const [isManual, setIsManual] = React.useState(false)
  const [manualAmounts, setManualAmounts] = React.useState<Record<string, string>>({})
  const queryClient = useQueryClient()

  const { data: members } = useQuery({
    queryKey: ["group", groupId, "members"],
    queryFn: () => api.get(`/groups/${groupId}/members`).then(res => res.data),
    enabled: !!groupId
  });

  // const members = React.useMemo(() => groupData?.members || [], [groupData?.members]);

  const form = useForm<ExpenseValues>({
    resolver: zodResolver(expenseSchema), defaultValues: {
      description: "",
      amount: "",
      payer_id: "",
      selectedMembers: [],
    },
  })

  React.useEffect(() => {
    if (isOpen) {
      if (isEdit && expense) {
        const selectedIds = expense.splits.map(s => s.user.id);

        form.reset({
          description: expense.description,
          amount: expense.amount.toString(),
          payer_id: expense.payer.id,
          selectedMembers: selectedIds
        })

        const initialManual: Record<string, string> = {}
        expense.splits.forEach(s => {
          initialManual[s.user.id] = s.amount.toString()
        })
        setManualAmounts(initialManual)
        setIsManual(true)
      } else {
        form.reset({
          description: "",
          amount: "",
          payer_id: members?.[0]?.id || "",
          selectedMembers: members?.map(m => m.id) || [],
        })
        setIsManual(false)
        setManualAmounts({})
      }
    }
  }, [isOpen, isEdit, expense, members, form])

  const watchedAmount = form.watch("amount")
  const watchedMembers = form.watch("selectedMembers")

  const mutation = useMutation({
    mutationFn: (data: CreateExpenseRequest) =>
      isEdit
        ? api.put(`/groups/${groupId}/expenses/${expense.id}`, data)
        : api.post(`/groups/${groupId}/expenses`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({queryKey: ["group", groupId]})
      if (isEdit) queryClient.invalidateQueries({queryKey: ["expense", groupId, expense.id]})

      toast.success(isEdit ? "Expense updated!" : "Expense added!")
      setIsOpen(false)
      if (!isEdit) form.reset()
    },
  })

  const autoSplits = React.useMemo(() => {
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

  const toggleMode = () => {
    if (!isManual) {
      const initialManual: Record<string, string> = {}
      autoSplits.forEach(s => initialManual[s.user_id] = s.amount)
      setManualAmounts(initialManual)
    }
    setIsManual(!isManual)
  }

  function onSubmit(values: ExpenseValues) {
    const finalSplits = isManual
      ? watchedMembers.map(id => ({user_id: id, amount: manualAmounts[id] || "0.00"}))
      : autoSplits

    if (isManual) {
      const sum = finalSplits.reduce((acc, curr) => acc + parseFloat(curr.amount), 0)
      if (Math.abs(sum - parseFloat(values.amount)) > 0.01) {
        toast.error(`Total sum (${sum.toFixed(2)}€) must match amount (${values.amount}€)`)
        return
      }
    }

    mutation.mutate({
      ...values,
      currency: "EUR",
      splits: finalSplits
    })
  }

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogContent className="sm:max-w-[450px] rounded-3xl overflow-hidden">
        <DialogHeader className="flex-row justify-start gap-2">
          <div
            className="h-12 w-12 rounded-2xl bg-pink-50 flex items-center justify-center text-pink-600 mb-2 border border-pink-100">
            <Receipt className="h-6 w-6"/>
          </div>
          <div className="flex-col">
            <DialogTitle className="text-2xl font-black">{isEdit ? "Edit Expense" : "Add Expense"}</DialogTitle>
            <DialogDescription>{isEdit ? "Modify expense details below." : "Split a new bill with your group."}</DialogDescription>
          </div>
        </DialogHeader>

        <form onSubmit={form.handleSubmit(onSubmit)}>
          <FieldGroup>
            <div className="grid grid-cols-3 gap-4">
              {/* Description */}
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

              {/* Amount */}
              <div className="col-span-1">
                <Controller
                  name="amount"
                  control={form.control}
                  render={({field, fieldState}) => (<Field data-invalid={fieldState.invalid}>
                    <FieldLabel>Amount</FieldLabel>
                    <div className="relative">
                      <Input {...field} placeholder="0.00" className="rounded-xl pr-7"/>
                      <span className="absolute right-3 top-2 text-sm">€</span>
                    </div>
                    {fieldState.invalid && <FieldError errors={[fieldState.error]}/>}
                  </Field>)}
                />
              </div>
            </div>

            {/* Payer Selection */}
            <Controller
              name="payer_id"
              control={form.control}
              render={({field}) => (
                <Field>
                  <FieldLabel>Paid by</FieldLabel>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select who paid"/>
                    </SelectTrigger>
                    <SelectContent>
                      {members && members.map((member) => (
                        <SelectItem
                          key={member.id}
                          value={member.id}
                        >
                          <div className="flex items-center gap-2">
                            <div
                              className="h-5 w-5 rounded-full bg-pink-100 !text-pink-600 flex items-center justify-center text-[8px] font-bold">
                              {member.name[0]}{member.surname[0]}
                            </div>
                            <span className="text-sm">{member.name} {member.surname}</span>
                          </div>
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </Field>
              )}
            />

            {/* Share Distribution */}
            <div>
              <div className="flex justify-between items-center">
                <FieldLabel>Share Distribution</FieldLabel>

                <div className="flex items-center gap-2 bg-zinc-50 px-3 py-2 rounded-2xl border border-zinc-100">
                  <span
                    className={`text-[10px] font-bold transition-colors ${!isManual ? "text-pink-600" : "text-zinc-400"}`}>
                    AUTO
                  </span>
                  <Switch
                    checked={isManual}
                    onCheckedChange={toggleMode}
                    className="data-[state=checked]:bg-pink-600"
                  />
                  <span
                    className={`text-[10px] font-bold transition-colors ${isManual ? "text-pink-600" : "text-zinc-400"}`}>
                    MANUAL
                  </span>
                </div>
              </div>

              <div className="max-h-[250px] overflow-y-auto">
                <UserListCard
                  users={members}
                  renderSubtext={(user) => {
                    const isSelected = watchedMembers.includes(user.id);
                    const amountNumber = parseFloat(watchedAmount);
                    const isAmountValid = !isNaN(amountNumber) && amountNumber > 0;

                    return (
                      <div className="h-8 flex items-center">
                        {!isSelected ? (
                          <span className="text-[11px] text-zinc-300 italic">Not included</span>
                        ) : isManual ? (
                          <div className="mt-1 relative w-24">
                            <Input
                              type="number"
                              step="0.01"
                              value={manualAmounts[user.id] || ""}
                              onChange={(e) => setManualAmounts(prev => ({...prev, [user.id]: e.target.value}))}
                              className="h-7 text-xs text-pink-600 pr-4 rounded-lg"
                            />
                            <span className="absolute right-2 top-1.5 text-pink-600">€</span>
                          </div>
                        ) : isAmountValid ? (
                          <span className="text-pink-600">
                          owes {autoSplits.find(s => s.user_id === user.id)?.amount} €
                        </span>
                        ) : null }
                      </div>
                    )
                  }}

                  renderActions={(user) => (
                    <Controller
                      name="selectedMembers"
                      control={form.control}
                      render={({field}) => {
                        const currentValues = Array.isArray(field.value) ? field.value : [];
                        return (
                          <Checkbox
                            checked={currentValues.includes(user.id)}
                            onCheckedChange={(checked) => {
                              const newValue = checked
                                ? [...currentValues, user.id]
                                : currentValues.filter(id => id !== user.id);
                              field.onChange(newValue);
                            }}
                          />
                        )
                      }}
                    />
                  )}
                />
              </div>

              {form.formState.errors.selectedMembers && (
                <p className="text-xs text-red-500 font-medium">{form.formState.errors.selectedMembers.message}</p>
              )}
            </div>
          </FieldGroup>

          <DialogFooter className="mt-4">
            <Button type="submit" variant="pinkPrimary" className="w-full" disabled={mutation.isPending}>
              {mutation.isPending ? "Adding..." : "Save Expense"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>)
}