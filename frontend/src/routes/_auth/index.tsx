import apiClient from "@/api/apiClient";
import { Button } from "@/components/ui/button";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/")({
  component: App,
});

function App() {
  async function Get() {
    const res = await apiClient.get("/");
    if (res) {
      console.log(res);
    } else {
      console.error("API is not reachable");
    }
  }

  return (
    <div className="text-center">
      <Button
        onClick={(e) => {
          e.preventDefault();
          Get();
        }}
      >
        Click this
      </Button>
    </div>
  );
}
