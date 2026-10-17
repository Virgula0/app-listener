# Graph Report - app-listener  (2026-10-17)

## Corpus Check
- 438 files · ~1,616,814 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 5569 nodes · 23617 edges · 187 communities (162 shown, 25 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2076 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `55e9ca98`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- findSectionStart
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- rcUser
- guardUnitTest
- globreserve_test.go
- go_pkg_path_filepath
- planUpdaters
- UI
- desktop.go
- os.FileMode
- orphans_test.go
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- NetGuardEvent
- guard.bpf.c
- binaryVetter
- Candidate
- inodeChain
- github.com/charmbracelet/bubbles/textarea.Model
- newTestTarget
- serve.go
- bpfstats/main.go
- .render
- picker
- github.com/spf13/cobra.Command
- IntegrationSuite
- Guard
- controlServer
- fileedit_model_test.go
- MountVouch
- IntegrationSuite
- IntegrationSuite
- Monitor
- IntegrationSuite
- update.go
- Ledger
- stdio
- audit.go
- catalogRefresher
- searchDirs
- GenerateConf
- .FilterExistingWhitelist
- daemonconfig.go
- github.com/charmbracelet/bubbletea.Cmd
- TrustGuard
- systemd.go
- go_pkg_os
- GuardInodeKey
- Vault
- buildOneGuard
- .Refresh
- exe_supersede.h
- buildTrustedSet
- installPage
- pinstate.go
- wm.go
- runMonitor
- changelog_test.go
- FileTarget
- fcntl
- copyRegularAt
- fileedit_symlink_test.go
- time.Time
- btrfsLayoutFrom
- sync.Mutex
- hashFileShared
- inode_dev
- NetEventType
- process_vm_readv.c
- GuardEvent
- check_watched
- time.Duration
- protected/daemon_test.go
- testKeyPair
- runDaemon
- Finder
- liveSession
- highlight_test.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- dialog_test.go
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- go_pkg_strings
- policyEditor
- NewDumpHook
- Well-known bypass classes
- appTheme
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- openBunTmp
- Build & sign release assets (reusable workflow)
- icon
- KernelDev
- mountsUnitTest
- NetworkMonitorRepository
- app-listener Terminal Demo Recording
- IntegrationSuite
- IntegrationSuite
- bufio.Reader
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- testLimiter
- New
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- User
- make test-integration (rootful Docker suite)
- IntegrationSuite
- supersede_probe.c
- newTestVetter
- trace-app-libs.sh
- IntegrationSuite
- IntegrationSuite
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- fakeGuardRepo
- io.Reader
- integrationtests/guard_test.go
- bpfLsmListed
- checkRunningBinaryMatchesInstalled
- net_tester/main.go
- inode_dev.h
- .Unlock
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- EventType
- layout.go
- IntegrationSuite
- Installation
- os.File
- graphify-refresh.sh
- Pulse
- parseGuardEvent
- TestMain
- WithSuspectsFrom
- feed.go
- shell
- IntegrationSuite
- Limitations
- Troubleshooting
- Daemon configuration
- edit-protected
- configedit_test.go
- Compatibility
- How it works
- SanitizeText
- bleedLayout
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- .runEditorHarness
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- launch_rule_at
- monitor: see what touches a path
- kernel-prof.sh
- Desktop GUI

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `Config` - 96 edges
3. `GuardInodeKey` - 96 edges
4. `Load()` - 90 edges
5. `IntegrationSuite` - 82 edges
6. `absPath()` - 81 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 68 edges
9. `shell` - 66 edges
10. `BinaryEntry` - 60 edges

## Surprising Connections (you probably didn't know these)
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `runConfigEditor()` --calls--> `EditText()`  [INFERRED]
  cmd/functions/install/steps.go → internal/install/editor.go
- `resolveAssets()` --calls--> `asset()`  [INFERRED]
  cmd/functions/update/update.go → internal/gui/icons.go

## Import Cycles
- None detected.

## Communities (187 total, 25 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.19
Nodes (26): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey, launchModel() (+18 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (123): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+115 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "findSectionStart"
Cohesion: 0.22
Nodes (9): findSectionStart(), isLibDirectiveText(), ParseSectionHeaderPath(), unquotePath(), IsLibraryDirective(), parseLibDirLine(), PruneLibraryBlock(), ParseSectionWhitelist() (+1 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.11
Nodes (34): code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image(), fill_inode_key(), get_task_exe_key() (+26 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "rcUser"
Cohesion: 0.12
Nodes (35): setupBunTmpdir(), revertBunLaunchers(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir() (+27 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (23): addErrKind, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest (+15 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.17
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "go_pkg_path_filepath"
Cohesion: 0.26
Nodes (13): applyLiveRefresh(), softenAutomatedRefreshErr(), TestApplyLiveRefreshDeliversReload(), TestSoftenAutomatedRefreshErr(), isAbsoluteCandidate(), TestCatalogSanity(), LiveEmptyWhitelistRejected(), TestLiveEmptyWhitelistRejected() (+5 more)

### Community 11 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 12 - "UI"
Cohesion: 0.08
Nodes (20): installUI, multiAnswer, cancelUI, scrollable(), ask(), diffText(), orDefault(), promptBody() (+12 more)

### Community 13 - "desktop.go"
Cohesion: 0.08
Nodes (26): installHooks(), releaseBuild(), Run(), stageItem, stageState, textAnswer, encodeTable(), invokingIDs() (+18 more)

### Community 14 - "os.FileMode"
Cohesion: 0.08
Nodes (34): TestDaemonUnitWithMetadataOutput(), dupFD(), installFile(), installFileAs(), installServices(), installSSHAgent(), installSSHAgentEnv(), installUnit() (+26 more)

### Community 15 - "orphans_test.go"
Cohesion: 0.21
Nodes (19): appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors(), orphanedPolicies() (+11 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.09
Nodes (5): IntegrationSuite, IntegrationSuite, guardDeltaEvents(), le64HexKey(), absPath()

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 19 - "fileEditModel"
Cohesion: 0.09
Nodes (14): TestRuntimeClassMarkers(), applyNewFileMeta(), clipLabel(), humanSize(), listEntries(), clamp(), sortNodes(), writeFileInPlace() (+6 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.09
Nodes (59): Vault, NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged() (+51 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (10): fakeDaemon, Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots() (+2 more)

### Community 23 - "auth.go"
Cohesion: 0.08
Nodes (47): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+39 more)

### Community 24 - "NetGuardEvent"
Cohesion: 0.09
Nodes (19): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), snapNetGuard, Mode (+11 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (53): add_inode_to_guard(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with(), event_is_read_class(), evict_inode_from_guard() (+45 more)

### Community 26 - "binaryVetter"
Cohesion: 0.05
Nodes (44): CheckBPFLSM(), newBinaryVetter(), openBinaryVetter(), backingDeviceUnion(), buildConcurrency(), buildGuards(), configureDaemonLogging(), confirmMasterKeyOverwrite() (+36 more)

### Community 27 - "Candidate"
Cohesion: 0.06
Nodes (51): systemPatterns(), askBunTmpdirUsers(), gatherPerUserSetup(), appendSectionsAndEdit(), applyDiffAdditions(), collectDiffAdditions(), configPaths(), pathCovered() (+43 more)

### Community 28 - "inodeChain"
Cohesion: 0.26
Nodes (13): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+5 more)

### Community 29 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.15
Nodes (19): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+11 more)

### Community 30 - "newTestTarget"
Cohesion: 0.09
Nodes (22): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+14 more)

### Community 31 - "serve.go"
Cohesion: 0.07
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 32 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 33 - ".render"
Cohesion: 0.14
Nodes (11): lineStyles, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle(), setGutter() (+3 more)

### Community 34 - "picker"
Cohesion: 0.09
Nodes (21): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, TestSearchTreeFollowsNoSymlink(), errText() (+13 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.11
Nodes (25): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+17 more)

### Community 37 - "Guard"
Cohesion: 0.05
Nodes (31): unlockUnderGuard(), vaultOpUnderGuard(), deferredBinary, rootHandle, BinaryEntry, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), checkBtrfsKeys() (+23 more)

### Community 38 - "controlServer"
Cohesion: 0.15
Nodes (10): swapFile(), readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), configEditor, controlServer (+2 more)

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.13
Nodes (31): defaultNewFileMode(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed() (+23 more)

### Community 40 - "MountVouch"
Cohesion: 0.09
Nodes (18): ensureBtrfsLayout(), deleteResKeys(), dropSupersededLocked(), PruneSuperseded(), supersedeMaps(), supersedePruneLoop(), warnLostStamps(), loadPinnedExeMaps() (+10 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 43 - "Monitor"
Cohesion: 0.06
Nodes (13): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), dirTarget(), TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace(), newMountTable() (+5 more)

### Community 44 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 45 - "update.go"
Cohesion: 0.05
Nodes (55): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), showChangelog(), assetsFor(), checkUpdatePreconditions(), compareStableVersions() (+47 more)

### Community 46 - "Ledger"
Cohesion: 0.09
Nodes (15): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), TestProgressReader(), fyneLogFilter (+7 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "audit.go"
Cohesion: 0.12
Nodes (26): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+18 more)

### Community 49 - "catalogRefresher"
Cohesion: 0.17
Nodes (7): drain(), newCatalogRefresher(), readEvents(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer

### Community 50 - "searchDirs"
Cohesion: 0.18
Nodes (10): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths(), resolveSoname() (+2 more)

### Community 51 - "GenerateConf"
Cohesion: 0.06
Nodes (47): TestDiffMergePreservesExistingAndParses(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGeneratedLibraryBlockRoundTrips() (+39 more)

### Community 52 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.13
Nodes (41): grantNewBinaries(), watchGroup, Lines(), TestLinesKeyOnLinkAndIncludePending(), addResource(), applyAllowLib(), applyConfigLine(), applyDirective() (+33 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (46): diffModel, editorModel, Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace() (+38 more)

### Community 55 - "TrustGuard"
Cohesion: 0.08
Nodes (13): exeRow, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent() (+5 more)

### Community 56 - "systemd.go"
Cohesion: 0.10
Nodes (33): buildBinaryIfNeeded(), ensureInstalledBinary(), ensureMasterKey(), installBinaryOnly(), prepareInstallation(), TestEnsureInstalledBinary(), mustCwd(), revertSystemFiles() (+25 more)

### Community 57 - "go_pkg_os"
Cohesion: 0.04
Nodes (50): Requested(), TestRequested(), TestValidateFlagsEditConfig(), TestPickEncryptedDirSingleSkipsPicker(), validateFlags(), pickEncryptedDir(), runOfflineEdit(), atomicWriteAt() (+42 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.04
Nodes (47): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+39 more)

### Community 59 - "Vault"
Cohesion: 0.15
Nodes (18): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq() (+10 more)

### Community 60 - "buildOneGuard"
Cohesion: 0.09
Nodes (32): multicallAdmissible(), buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards() (+24 more)

### Community 61 - ".Refresh"
Cohesion: 0.05
Nodes (10): clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, asset(), coloredIcon() (+2 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.17
Nodes (18): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+10 more)

### Community 63 - "buildTrustedSet"
Cohesion: 0.08
Nodes (24): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+16 more)

### Community 64 - "installPage"
Cohesion: 0.07
Nodes (15): binSource, installAction, installedBinary, installPage, logHook, logLine, runView, hashReader() (+7 more)

### Community 65 - "pinstate.go"
Cohesion: 0.17
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 67 - "runMonitor"
Cohesion: 0.07
Nodes (30): BuildEBPFTargets(), CheckEBPF(), RawTarget, isSubDir(), MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), pathWithinDir() (+22 more)

### Community 68 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 69 - "FileTarget"
Cohesion: 0.13
Nodes (12): Target, specFromForm(), specSummary(), EventTypes(), Files, FileSpec, FileTarget, Policy (+4 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "fileedit_symlink_test.go"
Cohesion: 0.20
Nodes (17): checkFreshDir(), openDirNoFollow(), renewDir(), openRootT(), TestWriteFileKeepMeta(), TestWriteFileKeepMetaRefusesSymlink(), mustRead(), mustWrite() (+9 more)

### Community 73 - "time.Time"
Cohesion: 0.08
Nodes (25): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+17 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.17
Nodes (9): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers() (+1 more)

### Community 75 - "sync.Mutex"
Cohesion: 0.07
Nodes (21): netGuardPill(), newActivity(), newExeResolver(), guardRow(), TestRingKeepsNewest(), clonePolicy(), NetEngines, NetExtra (+13 more)

### Community 76 - "hashFileShared"
Cohesion: 0.22
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 77 - "inode_dev"
Cohesion: 0.18
Nodes (15): exe_is_superseded(), exe_refused(), inode_dev(), check_and_emit_args(), code_suspect(), current_is_inspector(), exe_tag_bits(), fill_path() (+7 more)

### Community 78 - "NetEventType"
Cohesion: 0.07
Nodes (21): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), Cstr(), FormatAddr() (+13 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.07
Nodes (17): holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), snapGuard, TestDispatchFullChannelNeverDropsDenials() (+9 more)

### Community 81 - "check_watched"
Cohesion: 0.24
Nodes (14): admit_refusal(), check_watched(), emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg() (+6 more)

### Community 82 - "time.Duration"
Cohesion: 0.17
Nodes (9): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), procStats, statsSample, cpuTime(), humanBytes(), rssBytes() (+1 more)

### Community 83 - "protected/daemon_test.go"
Cohesion: 0.12
Nodes (26): deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall(), removeKeyAndEmptyDir(), removeMasterKey(), TestRemoveKeyAndEmptyDir(), Delete() (+18 more)

### Community 84 - "testKeyPair"
Cohesion: 0.23
Nodes (10): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+2 more)

### Community 85 - "runDaemon"
Cohesion: 0.09
Nodes (24): newControlManager(), startControlManager(), makeReloadHandler(), markReady(), noteReloadedSettings(), notifySystemdReady(), runDaemon(), startCatalogRefresh() (+16 more)

### Community 86 - "Finder"
Cohesion: 0.15
Nodes (19): findButton, span, findAll(), Finder, lowerRunes(), NewFinder(), overlay(), runesEqual() (+11 more)

### Community 87 - "liveSession"
Cohesion: 0.19
Nodes (11): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), displayEntries(), eventsLabel(), forwardBinaries(), holdForward() (+3 more)

### Community 88 - "highlight_test.go"
Cohesion: 0.13
Nodes (25): New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines(), TestSetTextStartsAtFirstLine() (+17 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "dialog_test.go"
Cohesion: 0.39
Nodes (10): openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder(), TestPickerSubfolderSearchPicksFile() (+2 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.09
Nodes (33): UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), catchLifecycleSignals(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+25 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.07
Nodes (19): foreignPinnedLink(), eventTypeKey(), mcNameKey(), statInodeKey(), codeVouched(), isTraced(), judgePasses(), loaderEnv() (+11 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "go_pkg_strings"
Cohesion: 0.11
Nodes (3): launchCase, repeatSpaces(), wrap()

### Community 100 - "policyEditor"
Cohesion: 0.13
Nodes (12): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+4 more)

### Community 101 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "icon"
Cohesion: 0.06
Nodes (36): binaryList, binaryRow, checkGrid, filesPage, handButton, homePage, iconTile, panel (+28 more)

### Community 110 - "KernelDev"
Cohesion: 0.07
Nodes (37): fdStat, BinaryStat, placedWalk, physicalParent(), statFD(), binaryStatOf(), isFuseType(), parseMajorMinor() (+29 more)

### Community 111 - "mountsUnitTest"
Cohesion: 0.10
Nodes (11): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+3 more)

### Community 112 - "NetworkMonitorRepository"
Cohesion: 0.12
Nodes (13): NetworkMonitorRepository, NewGuardUseCase(), NetworkMonitorUseCase, NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle() (+5 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 116 - "bufio.Reader"
Cohesion: 0.20
Nodes (11): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest(), TestReadGrantRequestForward() (+3 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 121 - "New"
Cohesion: 0.06
Nodes (53): harnessSuite, migrationTarget, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord() (+45 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "User"
Cohesion: 0.07
Nodes (36): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), systemRule(), entryWriters(), reserveChain(), symlinkTargetGlob(), underResource() (+28 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "newTestVetter"
Cohesion: 0.19
Nodes (18): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+10 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "fakeGuardRepo"
Cohesion: 0.07
Nodes (17): SupersededBinary, TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess (+9 more)

### Community 143 - "io.Reader"
Cohesion: 0.40
Nodes (4): entryOf(), EntryWithHash(), mountinfoDev(), TestMountinfoDev()

### Community 144 - "integrationtests/guard_test.go"
Cohesion: 0.36
Nodes (6): attrOpCase, eventsForPath(), guardEventTypes(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 145 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 146 - "checkRunningBinaryMatchesInstalled"
Cohesion: 0.40
Nodes (5): checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), fileSHA256(), sameContent(), TestCheckRunningBinaryMatchesInstalled()

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 149 - "inode_dev.h"
Cohesion: 0.24
Nodes (8): ctx_ptr(), sb_dev(), btrfs_copy_gate(), guard_file_ioctl(), guard_file_ioctl_compat(), guard_inode_free(), is_btrfs_copy_ioctl(), trust_inode_free()

### Community 150 - ".Unlock"
Cohesion: 0.11
Nodes (17): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn() (+9 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.29
Nodes (19): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+11 more)

### Community 154 - "EventType"
Cohesion: 0.10
Nodes (12): snapEngines, snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting() (+4 more)

### Community 155 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.05
Nodes (60): confirmUserPlaced(), VettedInode, Multicall, MulticallKind, WithVettedKeys(), ExeKey(), IsMulticallRefusal(), multicallFamily() (+52 more)

### Community 161 - "Pulse"
Cohesion: 0.25
Nodes (7): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub

### Community 162 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 163 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 164 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 165 - "feed.go"
Cohesion: 0.07
Nodes (38): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+30 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (33): InstallMode, InstallRequest, Options, pageDef, pageID, resizeWatch, shell, sizedDialog (+25 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 179 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 183 - "SanitizeText"
Cohesion: 0.09
Nodes (13): ownerPaths(), drainEphemeral(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), runHeadless(), runHeadless(), inspectorAdmitter, vetDecision (+5 more)

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 194 - "launch_rule_at"
Cohesion: 0.50
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 206 - "Desktop GUI"
Cohesion: 0.20
Nodes (10): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Tables (+2 more)

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 523 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **25 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Process inspectors` connect `time.Time` to `Daemon configuration`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.023622047244094488 - nodes in this community are weakly interconnected._
- **Why does `Daemon configuration` connect `Daemon configuration` to `time.Time`, `README.md`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._
- **Why does `watchPattern` connect `time.Time` to `globreserve_test.go`, `go_pkg_path_filepath`, `catalogRefresher`?**
  _High betweenness centrality (0.047) - this node is a cross-community bridge._
- **Should `daemonconfig_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07737874861162532 - nodes in this community are weakly interconnected._