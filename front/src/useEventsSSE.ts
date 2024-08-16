import {useEffect} from 'react';
import {sseTesting, EntityType} from "./api.ts";
import {useAppDispatch} from "./store.ts";

export const useSseForDocumentUpdates = (start: boolean, handler: () => void) => {
	const dispatch = useAppDispatch();
	useEffect(() => {
		if (!start) return;
		// console.log(start, 'startSse')
			const eventSource = new EventSource(`http://localhost:8080/events`);
			eventSource.onmessage = (ev: MessageEvent<EntityType[]>) => {
				console.log(ev.data[0].id)
				if (ev.data) {
					const parsedData:EntityType[] = JSON.parse(ev.data as unknown as string)
					console.log(parsedData, 'parsedDta')
					dispatch(
						sseTesting.util.updateQueryData('getConnect', undefined, (data) => {
							if (!Array.isArray(data)) {
								return;
							}
							if (data) {
								// console.log(data, 'data -> useSseForDocumentUpdates -> updateQueryData 2')
								const index = data.findIndex((el) => el.id === parsedData[0].id)
								console.log(index, 'index useSseForDocumentUpdates')
								data[index].status = parsedData[0].status
								// console.log(data, 'data -> useSseForDocumentUpdates')
							}
						})
					)

					eventSource.close()
				}
			}
			handler()

	}, [dispatch, handler, start]);

};
