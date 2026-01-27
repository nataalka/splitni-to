import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom"
import LoginPage from "@/pages/LoginPage"
import { Toaster } from "@/components/ui/sonner";
import RegisterPage from "@/pages/RegisterPage";
import GroupsPage from "@/pages/GroupsPage";
import GroupDetailPage from "@/pages/GroupDetailPage.tsx";
import FriendsPage from "@/pages/FriendsPage.tsx";
import MainLayout from "@/pages/MainLayout.tsx";
import ExpenseDetailPage from "@/pages/ExpenseDetailPage.tsx";

export function App() {
  return (<BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/login" replace/>}/>

        <Route path="/login" element={<LoginPage/>}/>

        <Route path="/register" element={<RegisterPage/>}/>

        <Route path="*" element={<div className="p-10 text-center">404 - Not Found</div>}/>

        <Route element={<MainLayout/>}>
          <Route path="/friends" element={<FriendsPage/>}/>

          <Route path="/groups" element={<GroupsPage/>}/>
          <Route path="/groups/:id" element={<GroupDetailPage/>}/>

          <Route path="/groups/:group_id/expenses/:expense_id" element={<ExpenseDetailPage/>}/>
        </Route>
      </Routes>
      <Toaster/>
    </BrowserRouter>)
}

export default App;