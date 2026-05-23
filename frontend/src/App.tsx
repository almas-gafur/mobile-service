import Dashboard from './pages/Dashboard';
import PublicApplication from './pages/PublicApplication';
import StatusPage from './pages/StatusPage';

export default function App() {
  if (window.location.pathname.startsWith('/track/')) {
    return <StatusPage />;
  }

  if (window.location.pathname.startsWith('/admin')) {
    return <Dashboard />;
  }

  return <PublicApplication />;
}
