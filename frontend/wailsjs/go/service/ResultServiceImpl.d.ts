// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {entities} from '../models';
import {context} from '../models';

export function GetAll(arg1:entities.RequestResultDTO):Promise<Array<entities.ResultDTO>>;

export function GetByPath(arg1:string):Promise<entities.ResultDTO>;

export function SetContext(arg1:context.Context):Promise<void>;
