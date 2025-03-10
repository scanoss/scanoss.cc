// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {context} from '../models';
import {service} from '../models';

export function BeforeClose(arg1:context.Context):Promise<boolean>;

export function GetRecentScanRoots():Promise<Array<string>>;

export function GetResultFilePath():Promise<string>;

export function GetScanRoot():Promise<string>;

export function GetScanSettingsFilePath():Promise<string>;

export function Init(arg1:context.Context,arg2:service.ScanossSettingsService,arg3:service.KeyboardService):Promise<void>;

export function JoinPaths(arg1:Array<string>):Promise<string>;

export function SelectDirectory():Promise<string>;

export function SelectFile(arg1:string):Promise<string>;

export function SetResultFilePath(arg1:string):Promise<void>;

export function SetScanRoot(arg1:string):Promise<void>;

export function SetScanSettingsFilePath(arg1:string):Promise<void>;
