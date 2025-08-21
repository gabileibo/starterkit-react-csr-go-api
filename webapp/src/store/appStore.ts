import { create } from 'zustand';

interface AppState {
  // Theme
  theme: 'light' | 'dark';
  toggleTheme: () => void;

  // Example state - replace with your own
  counter: number;
  increment: () => void;
  decrement: () => void;
}

// Initialize theme from localStorage or system preference
const getInitialTheme = (): 'light' | 'dark' => {
  if (typeof window === 'undefined') return 'light';

  const savedTheme = localStorage.getItem('theme') as 'light' | 'dark' | null;
  if (savedTheme) return savedTheme;

  return window.matchMedia('(prefers-color-scheme: dark)').matches
    ? 'dark'
    : 'light';
};

export const useAppStore = create<AppState>((set, get) => {
  const initialTheme = getInitialTheme();

  // Apply initial theme to document
  if (initialTheme === 'dark') {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }

  return {
    theme: initialTheme,
    toggleTheme: () => {
      const newTheme = get().theme === 'light' ? 'dark' : 'light';
      set({ theme: newTheme });

      // Apply theme to document
      if (newTheme === 'dark') {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }

      // Save to localStorage
      localStorage.setItem('theme', newTheme);
    },

    // Example state - replace with your own
    counter: 0,
    increment: () => set((state) => ({ counter: state.counter + 1 })),
    decrement: () => set((state) => ({ counter: state.counter - 1 })),
  };
});
