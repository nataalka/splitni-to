import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import LoginPage from "@/pages/LoginPage"
import { Toaster } from "@/components/ui/sonner";
import RegisterPage from "@/pages/RegisterPage";
import GroupsPage from "@/pages/GroupsPage";
import GroupDetailPage from "@/pages/GroupDetail";
import FriendsPage from "@/pages/FriendsPage.tsx";

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/login" replace/>}/>

        <Route path="/login" element={<LoginPage/>}/>

        <Route path="/register" element={<RegisterPage/>}/>

        <Route path="/friends" element={<FriendsPage/>}/>

        <Route path="/groups" element={<GroupsPage/>}/>
        <Route path="/groups/:id" element={<GroupDetailPage/>}/>

        <Route path="*" element={<div className="p-10 text-center">404 - Not Found</div>}/>
      </Routes>
      <Toaster/>
    </BrowserRouter>
  )
}

export default App;