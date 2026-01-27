import { format } from "date-fns"
import type { Expense } from "@/types"
import { Link } from "react-router-dom";

interface ExpenseListItemProps {
  expense: Expense;
  payerName: string;
  currentUserId: string | undefined;
}

export function ExpenseListItem({ expense, payerName, currentUserId }: ExpenseListItemProps) {
  const usersSplit = expense.splits?.find(s => s.user_id === currentUserId);
  const usersAmount = usersSplit ? Number(usersSplit.amount) : 0;
  const isPayer = expense.payer_id === currentUserId;

  const displayAmount = isPayer
    ? Number(expense.amount) - usersAmount
    : usersAmount;

  const notInvolved = !isPayer && !usersSplit;

  return (
    <Link
      to={`/groups/${expense.group_id}/expenses/${expense.id}`}
      className="block hover:opacity-80 transition-opacity"
    >
      <div className="group flex items-center justify-between p-4 hover:bg-pink-50/20 transition-all cursor-pointer">
        <div className="flex items-center gap-4">
          {/* Date */}
          <div
            className="h-12 w-12 rounded-xl bg-zinc-100 flex flex-col items-center justify-center text-zinc-500 border border-zinc-100 group-hover:border-pink-200 group-hover:bg-white transition-all">
          <span className="text-[9px] font-black uppercase leading-none tracking-tighter">
            {format(new Date(expense.created_at), "MMM")}
          </span>
            <span className="text-lg font-black leading-none group-hover:text-zinc-900">
            {format(new Date(expense.created_at), "dd")}
          </span>
          </div>

          {/* Description */}
          <div>
            <h3 className="font-bold text-zinc-900 leading-tight truncate">
              {expense.description}
            </h3>
            <p className="text-[11px] text-zinc-500 mt-1">
              Paid by <span className={`font-bold ${isPayer ? "text-pink-600" : "text-zinc-700"}`}>{payerName}</span>
            </p>
          </div>
        </div>

        <div className="flex gap-6 items-center">
          {/* Total amount */}
          <div className="text-right hidden sm:block">
            <p className="text-[9px] text-zinc-400 uppercase tracking-widest font-bold">
              Total bill
            </p>
            <p className="text-md font-bold text-zinc-500">
              {Number(expense.amount).toFixed(2)} €
            </p>
          </div>

          {/* Users share */}
          <div className="text-right min-w-[85px]">
            {notInvolved ? (
              <>
                <p className="text-[9px] uppercase tracking-widest font-bold mt-1 text-zinc-300">No debt</p>
                <p className="text-sm font-bold text-zinc-300 italic">not involved</p>
              </>
            ) : (
              <>
                <p
                  className={`text-[9px] uppercase tracking-widest font-bold ${isPayer ? "text-emerald-600/70" : "text-pink-600/70"}`}>
                  {isPayer ? "You lent" : "You owe"}
                </p>
                <p className={`text-lg font-bold leading-none ${isPayer ? "text-emerald-600" : "text-pink-600"}`}>
                  {isPayer ? "+" : "-"}{displayAmount.toFixed(2)} €
                </p>
              </>
            )}
          </div>
        </div>
      </div>
    </Link>
  )
}