import { AvatarImage } from "@radix-ui/react-avatar";
import type { GetShoppingOrderResponseWritable } from "@/api";
import {
	AvatarGroup,
	AvatarGroupTooltip,
} from "@/components/animate-ui/components/animate/avatar-group";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
	Card,
	CardContent,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import {
	Item,
	ItemContent,
	ItemDescription,
	ItemMedia,
	ItemTitle,
} from "@/components/ui/item";

type Props = {
	shoppingOrder: GetShoppingOrderResponseWritable;
};

const ShoppingCard = ({ shoppingOrder }: Props) => {
	return (
		<div className="flex flex-col gap-3">
			<h3 className="text-2xl font-bold">Shopping List</h3>
			<Card>
				<CardHeader>
					<CardTitle className="flex justify-between text-xl">
						<span>Total items: {shoppingOrder.items?.length ?? 0}</span>
					</CardTitle>
				</CardHeader>
				<CardContent className="flex flex-col gap-3 p-4">
					{shoppingOrder.items.map((item) => (
						<Item key={item.food_id} variant={"outline"}>
							<ItemMedia>
								<AvatarGroup>
									{item.users.map((user) => (
										<Avatar key={user.id}>
											<AvatarImage src={user.avatar_url} />
											<AvatarFallback>
												{user.display_name.slice(0, 2)}
											</AvatarFallback>
											<AvatarGroupTooltip>
												{user.display_name}
											</AvatarGroupTooltip>
										</Avatar>
									))}
								</AvatarGroup>
							</ItemMedia>
							<ItemContent>
								<ItemTitle>{item.food_name}</ItemTitle>
								<ItemDescription>Quantity: {item.quantity}</ItemDescription>
							</ItemContent>
						</Item>
					))}
				</CardContent>

				<CardFooter className="flex justify-end items-center">
					<div className="flex flex-col sm:flex-row items-center justify-end w-full text-lg">
						<span className="font-bold mr-2">Total Cost:</span>
						<span className="font-medium text-slate-700">
							{Intl.NumberFormat("vi-VN", {
								style: "currency",
								currency: "VND",
							}).format(shoppingOrder.total_price)}
						</span>
					</div>
				</CardFooter>
			</Card>
		</div>
	);
};

export default ShoppingCard;
