import * as React from "react";
import {
  Dialog,
  DialogContent, DialogDescription, DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import { GroupListCard } from "@/components/GroupListCard.tsx";

interface GroupSelectDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onGroupSelected: (groupId: string) => void;
}

export function GroupSelectDialog({ open, onOpenChange, onGroupSelected }: GroupSelectDialogProps) {
  const { data: groups } = useQuery({
    queryKey: ["groups"],
    queryFn: () => api.get("/groups").then((res) => res.data),
    enabled: open,
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[450px] rounded-3xl overflow-hidden">
        <DialogHeader className="justify-start gap-2">
          <DialogTitle className="text-2xl font-black">
            {"Select a Group"}
          </DialogTitle>
          <DialogDescription>{"Record new expense into a group."}</DialogDescription>
        </DialogHeader>
        <GroupListCard
          groups={groups}
          onGroupClick={(id) => onGroupSelected(id)}
          as="button"/>
      </DialogContent>
    </Dialog>
  );
}