import {useEffect} from 'react';
import {sseTesting, EntityType} from "../api/api.ts";
import {useAppDispatch} from "../common/store.ts";

type ResponseEntityType = { "processing_end": boolean } & EntityType

export const useSseForDocumentUpdates = (start: boolean, handler: () => void) => {
	const dispatch = useAppDispatch();
	useEffect(() => {
		if (!start) return;
		// console.log(start, 'startSse')
		const eventSource = new EventSource(`http://localhost:8080/events`);
		eventSource.onmessage = (ev: MessageEvent<ResponseEntityType>) => {
			console.log(ev.data.id)
			if (ev.data) {
				const parsedData: ResponseEntityType = JSON.parse(ev.data as unknown as string)

				console.log(parsedData, 'parsedDta')
				dispatch(
					sseTesting.util.updateQueryData('getConnect', undefined, (data) => {
						if (!Array.isArray(data)) {
							return;
						}
						if (data) {
							console.log(data, 'data -> useSseForDocumentUpdates -> updateQueryData 2')
							const index = data.findIndex((el) => el.id === parsedData.id)
							console.log(index, 'index useSseForDocumentUpdates')
							data[index].status = parsedData.status

							if (parsedData.processing_end) {
								eventSource.close()
								handler()
							}
							// console.log(data, 'data -> useSseForDocumentUpdates')
						}
					})
				)

				// eventSource.close()
			}
		}

		// eventSource.onopen = (ev: Event) => {
		// 	console.log(ev, 'eventSource.onopen => ev')
		// }

		// handler()

	}, [dispatch, handler, start]);

};
