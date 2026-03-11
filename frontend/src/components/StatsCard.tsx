import type { LucideIcon } from "lucide-react"
import { Card } from "@/components/ui/card"

interface StatsCardProps {
  label: string
  value: string | number
  icon: LucideIcon
  variant?: "default" | "positive" | "negative"
}

export function StatsCard({label, value, icon: Icon, variant = "default"}: StatsCardProps) {
  const styles = {
    positive: "bg-emerald-50/50 border-emerald-100 text-emerald-600 border-b-emerald-200",
    negative: "bg-pink-50/50 border-pink-100 text-pink-600 border-b-pink-200",
    default: "bg-white border-zinc-100 text-zinc-600 border-b-zinc-200",
  }[variant]

  const iconStyles = {
    positive: "bg-white border-emerald-100",
    negative: "bg-white border-pink-100",
    default: "bg-zinc-50 border-zinc-100",
  }[variant]

  return (
    <Card
      className={`rounded-2xl shadow-sm p-4 flex flex-col justify-between h-24 gap-0 border-b-2 transition-all ${styles}`}>
      <div className="flex items-end justify-between">
        <div className={`h-8 w-8 rounded-xl flex items-center justify-center border transition-colors ${iconStyles}`}>
          <Icon className="h-4 w-4"/>
        </div>
        <p className="text-[10px] font-extrabold uppercase tracking-widest">
          {label}
        </p>
      </div>

      <div className="text-right">
        <p className="text-2xl font-extrabold leading-none">
          {typeof value === "number" ? value.toFixed(2) : value} €
        </p>
      </div>
    </Card>
  )
}