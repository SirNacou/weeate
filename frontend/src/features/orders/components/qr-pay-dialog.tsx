import { Button } from "@/components/animate-ui/components/buttons/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTrigger,
} from "@/components/animate-ui/components/radix/dialog";
import { ImageWithFallback } from "@/components/image-with-fallback";

const QRPayDialog = () => {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button>Pay now</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <h3 className="text-lg font-semibold">Scan to Pay</h3>
        </DialogHeader>
        <div>
          <ImageWithFallback src={""} alt={"QR Code"}></ImageWithFallback>
        </div>
        <DialogFooter>
          <DialogClose>
            <Button variant="outline">Close</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default QRPayDialog;