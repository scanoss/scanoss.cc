// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {entities} from '../models';
import {adapter} from '../models';
import {common} from '../models';

export function ComponentGet(arg1:string):Promise<entities.ComponentDTO>;

export function FileGetLocalContent(arg1:string):Promise<adapter.FileDTO>;

export function FileGetRemoteContent(arg1:string):Promise<adapter.FileDTO>;

export function GetFilesToBeCommited():Promise<Array<adapter.GitFileDTO>>;

export function Greet(arg1:string):Promise<string>;

export function ResultGetAll(arg1:common.RequestResultDTO):Promise<Array<common.ResultDTO>>;
