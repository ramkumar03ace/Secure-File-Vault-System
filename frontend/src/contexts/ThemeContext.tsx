// src/contexts/ThemeContext.tsx

import React, {
  createContext,
  useState,
  useMemo,
  useContext,
  useEffect,
} from "react";

type Theme = "dark" | "light";

interface ThemeContextType {
  theme: Theme;
  toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export const ThemeProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  // 1. Initialize state from localStorage or default to 'dark'
  const [theme, setTheme] = useState<Theme>(
    () => (localStorage.getItem("theme") as Theme) || "dark"
  );

  // 2. The useEffect hook to apply the theme to the <html> element
  useEffect(() => {
    const root = window.document.documentElement;

    // Remove the other class to avoid conflicts
    root.classList.remove(theme === "light" ? "dark" : "light");
    // Add the current theme's class
    root.classList.add(theme);

    // 3. Save the user's preference to localStorage
    localStorage.setItem("theme", theme);
  }, [theme]);

  const toggleTheme = () => {
    setTheme((prevTheme) => (prevTheme === "light" ? "dark" : "light"));
  };

  const value = useMemo(() => ({ theme, toggleTheme }), [theme]);

  return (
    <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>
  );
};

export const useTheme = () => {
  const context = useContext(ThemeContext);
  if (context === undefined) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }
  return context;
};