import { %sColumns } from '@/columns/%s';
import ShowDetailLayout from "@/layouts/show/DetailLayout";
import { get%sId } from "@/services/meta/%s";

const DetailPage: React.FC<API.%s> = () => {
    return <ShowDetailLayout title="%s Detail" func={get%sId} columns={%sColumns}/>
};
export default DetailPage;
