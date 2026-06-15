import { useState, useEffect } from "react";
import { Navigate } from "react-router-dom";
import Sidebar from "../components/layout/Sidebar";
import DashboardContent from "./DashboardPage";
import { useAuthStore } from "../hooks/useAuthStore";
import type { FileCategory } from "../types";

export default function AppLayout() {
  const [activeCategory, setActiveCategory] = useState<FileCategory>("all");
  const { isAuthenticated } = useAuthStore();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return (
    <div className="h-screen flex overflow-hidden bg-surface-950">
      <Sidebar
        activeCategory={activeCategory}
        onCategoryChange={setActiveCategory}
      />
      <DashboardContent activeCategory={activeCategory} />
    </div>
  );
}
