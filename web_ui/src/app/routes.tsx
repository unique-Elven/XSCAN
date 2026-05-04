import { createBrowserRouter } from "react-router";
import { Layout } from "./components/Layout";
import { Scanner } from "./pages/Scanner";
import { Settings } from "./pages/Settings";
import { History } from "./pages/History";
import { ScannerProvider } from "./context/ScannerContext";

export const router = createBrowserRouter([
  {
    path: "/",
    element: (
      <ScannerProvider>
        <Layout />
      </ScannerProvider>
    ),
    children: [
      { index: true, element: <Scanner /> },
      { path: "settings", element: <Settings /> },
      { path: "history", element: <History /> },
    ],
  },
]);
