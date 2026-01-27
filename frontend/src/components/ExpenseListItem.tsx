import { format } from "date-fns"
import { Receipt, Calendar } from "lucide-react"
import type { Expense } from "@/types"

interface ExpenseListItemProps {
  expense: Expense;
}

export function ExpenseListItem({ expense }: ExpenseListItemProps) {
  return (
    <div className="group flex items-center justify-between p-4 hover:bg-zinc-50 transition-colors cursor-pointer">
      <div className="flex items-center gap-4">
        {/* Ikona účtenky */}
        <div className="h-12 w-12 rounded-2xl bg-zinc-50 flex flex-col items-center justify-center text-zinc-400 border border-zinc-100 group-hover:bg-white group-hover:text-pink-600 transition-all">
          <span className="text-[10px] font-bold uppercase leading-none">
            {format(new Date(expense.created_at), "MMM")}
          </span>
          <span className="text-lg font-black leading-none">
            {format(new Date(expense.created_at), "dd")}
          </span>
        </div>

        <div className="min-w-0">
          <h3 className="font-bold text-zinc-900 leading-tight truncate">
            {expense.description}
          </h3>
          <p className="text-xs text-zinc-500 mt-1">
            Paid by <span className="font-medium text-zinc-700">You</span> {/* Tu neskôr prepojíme meno platiteľa */}
          </p>
        </div>
      </div>

      <div className="text-right">
        <p className="font-black text-zinc-900 text-lg">
          {Number(expense.amount).toFixed(2)} {expense.currency === "EUR" ? "€" : expense.currency}
        </p>
        <p className="text-[10px] text-zinc-400 uppercase tracking-widest font-bold">
          Total bill
        </p>
      </div>
    </div>
  )
}