import { BrowserRouter, Routes, Route } from "react-router-dom";
import Profile from "./pages/Profile";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/:username" element={<Profile />} />
      </Routes>
    </BrowserRouter>
  );
}