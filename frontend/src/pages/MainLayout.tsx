import { Outlet, useNavigate, useLocation } from "react-router-dom";
import { ChevronLeft, Users, Plus, LayoutGrid } from "lucide-react";
import { Button } from "@/components/ui/button";
import { GroupSelectDialog } from "@/dialogs/GroupSelectDialog.tsx";
import { ExpenseFormDialog } from "@/dialogs/ExpenseFormDialog.tsx";
import * as React from "react";

export default function MainLayout() {
  const navigate = useNavigate();
  const location = useLocation();

  const canGoBack = location.pathname !== "/groups" && location.pathname !== "/friends";
  const [isGroupSelectOpen, setIsGroupSelectOpen] = React.useState(false);
  const [selectedGroupId, setSelectedGroupId] = React.useState<string | null>(null);
  const [isExpenseOpen, setIsExpenseOpen] = React.useState(false);

  return (
    <div className="flex flex-col min-h-screen bg-zinc-50 font-sans">
      {/* Fixed Header */}
      <header className="fixed top-0 left-0 right-0 h-16 bg-white border-b flex items-center px-4 z-50">
        <div className="flex-1">
          {canGoBack && (
            <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
              <ChevronLeft className="h-6 w-6 text-zinc-600" />
            </Button>
          )}
        </div>

        <div className="flex-1 text-center">
          <span className="font-black text-xl tracking-tighter text-pink-600">
            splitni.to
          </span>
        </div>

        <div className="flex-1 flex justify-end">
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 pt-20 pb-24 px-4 max-w-2xl mx-auto w-full">
        <Outlet />
      </main>

      {/* Bottom Navigation */}
      <nav className="fixed bottom-0 left-0 right-0 h-16 bg-white border-t flex items-center justify-around px-6 z-50 pb-safe">
        <button
          onClick={() => navigate("/friends")}
          className={`flex flex-col items-center gap-1 ${location.pathname === "/friends" ? "text-pink-600" : "text-zinc-400"}`}
        >
          <Users className="h-6 w-6" />
          <span className="text-[10px] font-medium">Friends</span>
        </button>

        <div className="relative -top-5">
          <button
            onClick={() => setIsGroupSelectOpen(true)}
            className="h-14 w-14 bg-pink-600 rounded-full flex items-center justify-center text-white shadow-lg shadow-pink-200 active:scale-95 transition-transform"
          >
            <Plus className="h-8 w-8" />
          </button>
        </div>

        <button
          onClick={() => navigate("/groups")}
          className={`flex flex-col items-center gap-1 ${location.pathname.includes("/groups") ? "text-pink-600" : "text-zinc-400"}`}
        >
          <LayoutGrid className="h-6 w-6" />
          <span className="text-[10px] font-medium">Groups</span>
        </button>
      </nav>

      <GroupSelectDialog
        open={isGroupSelectOpen}
        onOpenChange={setIsGroupSelectOpen}
        onGroupSelected={(groupId) => {
          setSelectedGroupId(groupId);
          setIsGroupSelectOpen(false);
          setIsExpenseOpen(true);
        }}
      />

      {selectedGroupId && (
        <ExpenseFormDialog
          groupId={selectedGroupId}
          open={isExpenseOpen}
          onOpenChange={setIsExpenseOpen}
        />
      )}
    </div>
  );
}