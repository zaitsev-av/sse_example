import {useEffect} from 'react';
import {sseTesting, EntityType} from "../api";
import {useAppDispatch} from "../common/store.ts";

type ResponseEntityType = { "processing_end": boolean } & EntityType

export const useSseForDocumentUpdates = (start: boolean, handler: () => void) => {
	const dispatch = useAppDispatch();
	useEffect(() => {
		if (!start) return;
		const eventSource = new EventSource(`http://localhost:8080/events`);
		eventSource.onmessage = (ev: MessageEvent<ResponseEntityType>) => {
			if (ev.data) {
				const parsedData: ResponseEntityType = JSON.parse(ev.data as unknown as string)

				dispatch(
					sseTesting.util.updateQueryData('getConnect', undefined, (data) => {
						if (!Array.isArray(data)) {
							return;
						}
						if (data) {
							const index = data.findIndex((el) => el.id === parsedData.id)
							data[index].status = parsedData.status

							if (parsedData.processing_end) {
								eventSource.close()
								handler()
							}
						}
					})
				)
			}
		}

	}, [dispatch, handler, start]);

};
