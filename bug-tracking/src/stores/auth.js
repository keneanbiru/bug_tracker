import { defineStore } from "pinia";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    user: null,
    token: localStorage.getItem("token") || null,
    role: localStorage.getItem("role") || null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.role === "admin",
    isDeveloper: (state) => state.role === "developer",
    isManager: (state) => state.role === "manager",
    currentUser: (state) => state.user,
    canReportBugs: (state) =>
      ["developer", "manager", "admin"].includes(state.role),
    canAssignBugs: (state) => ["manager", "admin"].includes(state.role),
    canEditBugStatus: (state) =>
      ["developer", "manager", "admin"].includes(state.role),
    canViewAllBugs: (state) => ["manager", "admin"].includes(state.role),
  },

  actions: {
    async login(credentials) {
      try {
        // TODO: Replace with actual API call
        const response = await fetch("http://localhost:8080/api/auth/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(credentials),
        });

        if (!response.ok) {
          throw new Error("Login failed");
        }

        const data = await response.json();
        this.token = data.token;
        this.user = data.user;
        this.role = data.user.role;

        localStorage.setItem("token", data.token);
        localStorage.setItem("role", data.user.role);

        return true;
      } catch (error) {
        console.error("Login error:", error);
        return false;
      }
    },

    async register(userData) {
      try {
        console.log("Attempting registration with:", userData);
        const response = await fetch(
          "http://localhost:8080/api/auth/register",
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(userData),
          }
        );

        if (!response.ok) {
          const errorData = await response.json().catch(() => null);
          console.error("Registration failed:", {
            status: response.status,
            statusText: response.statusText,
            error: errorData,
          });
          throw new Error(errorData?.message || "Registration failed");
        }

        const data = await response.json();
        console.log("Registration successful:", data);
        return { success: true, data };
      } catch (error) {
        console.error("Registration error:", error);
        return { success: false, error: error.message };
      }
    },

    logout() {
      this.user = null;
      this.token = null;
      this.role = null;
      localStorage.removeItem("token");
      localStorage.removeItem("role");
    },
  },
});
