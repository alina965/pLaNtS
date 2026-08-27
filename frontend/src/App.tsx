import {
  Navigate,
  Route,
  Routes,
} from "react-router";

import RegisterPage from "./page/RegisterPage";
import DashboardPage from "./page/DashboardPage";

function App() {
  return (
    <Routes>
      <Route path="/auth" element={<RegisterPage />} />
      <Route path="/dashboard" element={<DashboardPage />} />
      <Route path="/*" element={<Navigate to="/auth" />} />
    </Routes>
  );
}

export default App;