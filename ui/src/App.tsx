import { Routes, Route, BrowserRouter, Outlet } from 'react-router';
import { LoginPage } from '@/pages/login';
import { HomePage } from '@/pages/home';
import { AuthenticatedLayout } from '@/layouts/authenticated';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ThemeProvider } from '@/components/theme-provider';
import { ModeToggle } from '@/components/mode-toggler';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 60 * 1000,
      gcTime: 10 * 60 * 1000
    }
  }
})

function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <ModeToggle />
          <Routes>
            {/* <Route path="/" element={<LoginPage />} />
          <Route element={<AuthenticatedLayout />}>
            <Route path="/home" element={<HomePage />} />
          </Route> */}
            <Route path="/" element={<HomePage />} />
          </Routes>
        </BrowserRouter>
      </QueryClientProvider>
    </ThemeProvider>
  )
}

export default App
