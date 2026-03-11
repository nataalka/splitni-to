import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import type { User } from "@/types";

export function useAuth() {
  const { data: user, isLoading } = useQuery<User>({
    queryKey: ["auth-user"],
    queryFn: () => api.get("/users/me").then(res => res.data),
    staleTime: Infinity,
  });

  return {
    user,
    userId: user?.id,
    isLoading,
    isAuthenticated: !!user
  };
}