import type { LucideIcon } from "lucide-react";
import { Card } from "@/components/ui/card.tsx";
import * as React from "react";

interface EmptyBoxProps {
  description: string;
  Icon?: LucideIcon;
}

export const EmptyBox = ({ description, Icon }: EmptyBoxProps) => {
  return (
    <Card>
      <div className="text-center">
        { Icon && (
          <div className="bg-zinc-50 h-16 w-16 rounded-full flex items-center justify-center mx-auto mb-4">
            <Icon className="h-8 w-8 text-zinc-200" />
          </div>
        )}
        <p className="text-zinc-400 text-sm italic px-6">
          {description}
        </p>
      </div>
    </Card>
  )
}