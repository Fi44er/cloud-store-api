import { create } from "zustand";
import type { User } from "../types";
import { authApi } from "../api";

interface AuthState {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;

  login: (email: string, password: string) => Promise<void>;
  register: (email: string, username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  loadUser: () => Promise<void>;
  setUser: (user: User) => void;
}

export const useAuthStore = create<AuthState>((set: any) => ({
  user: null,
  isLoading: false,
  isAuthenticated: false,

  login: async (email, password) => {
    set({ isLoading: true });
    try {
      const { flow_id } = await authApi.initLoginFlow();
      const res = await authApi.login(flow_id, email, password);
      if (res && res.user) {
        // После логина Kratos ставит куку, мы просто загружаем пользователя
        const userStore = useAuthStore.getState();
        await userStore.loadUser();
      }
    } finally {
      set({ isLoading: false });
    }
  },

  register: async (email, username, password) => {
    set({ isLoading: true });
    try {
      const { flow_id } = await authApi.initRegistrationFlow();
      await authApi.register(flow_id, email, username, password);
      // После регистрации может потребоваться верификация или автоматический логин
      // Для простоты пробуем загрузить пользователя
      const userStore = useAuthStore.getState();
      await userStore.loadUser();
    } finally {
      set({ isLoading: false });
    }
  },

  logout: async () => {
    try {
      await authApi.logout();
    } finally {
      set({ user: null, isAuthenticated: false });
      window.location.href = "/login";
    }
  },

  loadUser: async () => {
    set({ isLoading: true });
    try {
      const session = await authApi.whoAmI();
      if (session && session.active) {
        const user: User = {
          id: session.identity.id,
          email: session.identity.traits.email || "",
          username: session.identity.traits.username || session.identity.traits.email || "Unknown",
          traits: session.identity.traits,
        };
        set({ user, isAuthenticated: true });
      } else {
        set({ user: null, isAuthenticated: false });
      }
    } catch (error) {
      set({ user: null, isAuthenticated: false });
    } finally {
      set({ isLoading: false });
    }
  },

  setUser: (user: User) => set({ user, isAuthenticated: !!user }),
}));
