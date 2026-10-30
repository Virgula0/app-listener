# Graph Report - app-listener  (2026-10-30)

## Corpus Check
- 562 files · ~1,788,665 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7308 nodes · 33396 edges · 234 communities (194 shown, 40 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2994 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4aa78abe`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- op
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- BinaryEntry
- mkdirs
- BusClient
- User
- protectedPage
- batchInodeRows
- Network
- NetworkMonitor
- IntegrationSuite
- guard_trust.bpf.c
- navItem
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- FileTarget
- guard.bpf.c
- LibraryClosure
- backend.go
- time.Time
- go_pkg_strings
- Binary
- NetGuardEvent
- fileedit_model_test.go
- Guard
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- controlServer
- os.FileMode
- .refreshOnce
- IntegrationSuite
- absPath
- monitorUnitTest
- fyne.io/fyne/v2/driver/desktop.Cursor
- SanitizeText
- TempBinary
- stdio
- update_test.go
- fadeScroll
- recorder/plan.go
- .showSized
- Client
- daemonconfig.go
- sync.Mutex
- .apply
- changelog.go
- mutedText
- GuardInodeKey
- Monitor
- GenerateConf
- .FilterExistingWhitelist
- inode_dev.h
- .filterEvent
- policyEditor
- unix_connect.c
- wm_test.go
- IntegrationSuite
- runForward
- runOfflineEdit
- Vault
- IntegrationSuite
- journal
- fcntl
- runMonitor
- Config
- codecheck/main.go
- runNetworkMonitor
- CandidateDir
- process_vm_readv.c
- os.File
- recordFlow
- newTestTarget
- rcUser
- golang.org/x/sys/unix.Stat_t
- bpfstats/main.go
- exe_supersede.h
- loginServer
- vaultfs_test.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- runView
- FetchLatest
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- canonicalBinaryPath
- root.go
- Well-known bypass classes
- binaryList
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- newTestVetter
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- highlight_test.go
- ConfinedEntry
- installPage
- Vault
- IntegrationSuite
- maintainerKeeper
- deploy
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- UI
- gui/session.go
- WebSocket /ws client
- inet_connect.c
- install.sh
- mc_tag.h
- classify
- make test-integration (rootful Docker suite)
- IntegrationSuite
- catalogRefresher
- logLine
- trace-app-libs.sh
- IntegrationSuite
- Finder
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- serveWebSocket
- serve.go
- liveSession
- github.com/charmbracelet/bubbles/textarea.Model
- Ledger
- appTheme
- GuardEvent
- GlobReservations
- find
- fakeStore
- bufio.Reader
- IntegrationSuite
- buildOneGuard
- watchSet
- copytree.go
- Daemon configuration
- Installation
- backend_test.go
- graphify-refresh.sh
- changelog_test.go
- NetEventType
- Gate
- debouncer
- feed.go
- shell
- layout.go
- procstats.go
- options.go
- SelectOrphans
- Limitations
- go_pkg_os
- IntegrationSuite
- TrustGuard
- net_tester/main.go
- NewEventFanout
- run
- runNetworkGuard
- AddServeFlags
- How it works
- inodeChain
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- github.com/spf13/cobra.Command
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- IntegrationSuite
- testLimiter
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- New
- pathCache
- pinstate_test.go
- DaemonEvent
- WithBottomBar
- EventType
- guardLSMHooks
- monitor
- IntegrationSuite
- logGrid
- TestReadFrame
- .eventTypeAllowed
- KernelDev
- VouchedDevs
- IntegrationSuite
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- SecretStoreNotes
- parseDaemonEvents
- saslServer
- downloadFile
- supersede_probe.c
- github.com/charmbracelet/bubbletea.Cmd
- TestCleanRel
- mountTable
- github.com/Virgula0/app-listener/tools/codecheck
- Troubleshooting
- Candidate
- IntegrationSuite
- taint_victim.c
- monitor: see what touches a path
- bleedLayout
- plan_test.go

## God Nodes (most connected - your core abstractions)
1. `Config` - 146 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 112 edges
4. `shell` - 97 edges
5. `Load()` - 95 edges
6. `absPath()` - 93 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 74 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `eBPF at three levels` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go
- `Why not an AppArmor profile?` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go
- `Distributions` --references--> `monitor()`  [INFERRED]
  docs/compatibility.md → integrationtests/secret_probe/main.go
- `Files` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `Workspaces` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go

## Import Cycles
- None detected.

## Communities (234 total, 40 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.19
Nodes (26): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+18 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (168): plantBun(), TestOpenBunTmp_KeepsVettedDir(), newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload() (+160 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): Hash(), configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+1 more)

### Community 3 - "op"
Cohesion: 0.12
Nodes (17): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+9 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.06
Nodes (92): TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+84 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.09
Nodes (51): prctl_gate(), guard_task_prctl(), admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader() (+43 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "BinaryEntry"
Cohesion: 0.05
Nodes (25): addErrKind, ReplacementCheck, BinaryEntry, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest (+17 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (35): TestOpenBunTmp_SymlinkReplacedTargetUntouched(), newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds() (+27 more)

### Community 10 - "BusClient"
Cohesion: 0.27
Nodes (7): BusClient, AttestBusClient(), busClientKey(), execStamp(), imageOf(), pidfdPid(), procIdentity()

### Community 11 - "User"
Cohesion: 0.06
Nodes (63): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), offerDaemonStart() (+55 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (24): backendLink, configView, forwardForm, forwardView, handCheck, protectedPage, recordForm, recordReview (+16 more)

### Community 13 - "batchInodeRows"
Cohesion: 0.24
Nodes (11): inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), walkInodeBatches(), countInodeRows() (+3 more)

### Community 14 - "Network"
Cohesion: 0.14
Nodes (7): netWarnings(), netGuardPill(), clonePolicy(), NetExtra, NetPolicy, Network, verifyHashes()

### Community 15 - "NetworkMonitor"
Cohesion: 0.15
Nodes (7): BinariesSummary(), TestBinariesSummary(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (64): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+56 more)

### Community 18 - "navItem"
Cohesion: 0.03
Nodes (21): binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel, panelRenderer (+13 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.11
Nodes (55): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+47 more)

### Community 22 - "Resource"
Cohesion: 0.05
Nodes (32): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+24 more)

### Community 23 - "auth_test.go"
Cohesion: 0.13
Nodes (38): CheckPassword(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf(), parseHash() (+30 more)

### Community 24 - "FileTarget"
Cohesion: 0.11
Nodes (13): guardPill(), Files, FileSpec, FileTarget, Policy, guardRow(), newFileTarget(), NormalizeSpec() (+5 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.07
Nodes (93): exe_refused(), inode_dev(), add_inode_to_guard(), admit_port(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect() (+85 more)

### Community 26 - "LibraryClosure"
Cohesion: 0.07
Nodes (27): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), IsScript(), ldConfInclude(), ldSoConfDirs() (+19 more)

### Community 27 - "backend.go"
Cohesion: 0.04
Nodes (98): confirmMasterKeyOverwrite(), runGenKey(), runOneShotMode(), installGUI(), asRoot(), buildBinaryIfNeeded(), certifyInstall(), checkInstalledMatchesSource() (+90 more)

### Community 28 - "time.Time"
Cohesion: 0.08
Nodes (13): LockedError, sweepID, sweepMemo, codeEntry, editView, Entry, fileSize(), joinRel() (+5 more)

### Community 29 - "go_pkg_strings"
Cohesion: 0.06
Nodes (16): parseReservedPorts(), readReservedPorts(), TestParseReservedPorts(), atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget() (+8 more)

### Community 30 - "Binary"
Cohesion: 0.08
Nodes (31): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), finishBinary(), Binary (+23 more)

### Community 31 - "NetGuardEvent"
Cohesion: 0.10
Nodes (17): snapNetGuard, NetGuardEvent, NetworkGuardRepository, NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo() (+9 more)

### Community 32 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 33 - "Guard"
Cohesion: 0.05
Nodes (24): deferredBinary, exeRow, rootHandle, supersededBinary, tempMaskOp, tempRow, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted() (+16 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (31): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, lstatOf(), outsideHardlinks() (+23 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.26
Nodes (23): busPath(), connexec(), connFD(), createItem(), dial(), execInto(), failf(), get() (+15 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (14): checkSecret(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp, backendState (+6 more)

### Community 38 - "controlServer"
Cohesion: 0.16
Nodes (7): authSelfPeer(), readClientLines(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), controlServer, editControlSession

### Community 39 - "os.FileMode"
Cohesion: 0.11
Nodes (26): TestDaemonUnitWithMetadataOutput(), installConfig(), installFileAs(), installUnit(), readAtNoFollow(), upsertFile(), upsertUnitAt(), writeAtNoFollow() (+18 more)

### Community 40 - ".refreshOnce"
Cohesion: 0.31
Nodes (8): logRefreshChange(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange(), SectionChange

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.05
Nodes (11): IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, launchCase, absPath(), TestMain() (+3 more)

### Community 43 - "monitorUnitTest"
Cohesion: 0.16
Nodes (4): sameInode(), dirTarget(), newMountTable(), monitorUnitTest

### Community 44 - "fyne.io/fyne/v2/driver/desktop.Cursor"
Cohesion: 0.12
Nodes (7): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), newButtonCursor(), TestButtonCursorFollowsState()

### Community 45 - "SanitizeText"
Cohesion: 0.07
Nodes (27): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), scriptIdentity(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome() (+19 more)

### Community 46 - "TempBinary"
Cohesion: 0.07
Nodes (21): fakeDaemon, binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), ResolveTempBinary(), stripPinnedAllow() (+13 more)

### Community 48 - "update_test.go"
Cohesion: 0.09
Nodes (20): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel() (+12 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (30): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+22 more)

### Community 50 - "recorder/plan.go"
Cohesion: 0.37
Nodes (4): recordTick, LiveEmptyWhitelistRejected(), TestLiveEmptyWhitelistRejected(), hint

### Community 51 - ".showSized"
Cohesion: 0.15
Nodes (12): installUI, multiAnswer, sizedDialog, scrollable(), selectable(), newFadeScroll(), ask(), promptBody() (+4 more)

### Community 52 - "Client"
Cohesion: 0.09
Nodes (14): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+6 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.08
Nodes (61): SecretScope, watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyBlockLine(), applyConfigLine(), applyDirective() (+53 more)

### Community 54 - "sync.Mutex"
Cohesion: 0.11
Nodes (10): fyneLogFilter, peerCred(), TestClientAuth(), TestSASL(), clientAuth(), readFrame(), serverAuth(), queued (+2 more)

### Community 55 - ".apply"
Cohesion: 0.17
Nodes (5): applier, ledgerOps, PlannedBinary, Store, vaultOps

### Community 56 - "changelog.go"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 57 - "mutedText"
Cohesion: 0.07
Nodes (47): checkGrid, filesPage, newLogo(), newBinaryList(), newBinaryRow(), newCheck(), newCheckGroup(), newRadioGroup() (+39 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (54): engine, portOwner, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel() (+46 more)

### Community 59 - "Monitor"
Cohesion: 0.15
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 60 - "GenerateConf"
Cohesion: 0.06
Nodes (42): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), discordCatalogEntry() (+34 more)

### Community 61 - ".FilterExistingWhitelist"
Cohesion: 0.12
Nodes (19): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+11 more)

### Community 62 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 64 - "policyEditor"
Cohesion: 0.13
Nodes (9): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), confirmWarnings(), sameEntries() (+1 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "wm_test.go"
Cohesion: 0.11
Nodes (24): reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket(), poll() (+16 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "runForward"
Cohesion: 0.10
Nodes (23): TestValidatePassword(), ValidatePassword(), dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig() (+15 more)

### Community 69 - "runOfflineEdit"
Cohesion: 0.09
Nodes (23): auditAfterEdit(), auditFindings(), reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive() (+15 more)

### Community 70 - "Vault"
Cohesion: 0.06
Nodes (36): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), decryptDirectories(), TestDecryptDirectoriesEmpty(), deprovisionKind, migrationTarget, isFileVaultCiphertext() (+28 more)

### Community 71 - "IntegrationSuite"
Cohesion: 0.16
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 72 - "journal"
Cohesion: 0.14
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "runMonitor"
Cohesion: 0.19
Nodes (14): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+6 more)

### Community 75 - "Config"
Cohesion: 0.04
Nodes (70): newCatalogRefresher(), sameConfig(), newControlManager(), startControlManager(), startControlServer(), makeReloadHandler(), startCatalogRefresh(), inspectorPaths() (+62 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "runNetworkMonitor"
Cohesion: 0.10
Nodes (13): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString() (+5 more)

### Community 78 - "CandidateDir"
Cohesion: 0.07
Nodes (27): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName() (+19 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.21
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "os.File"
Cohesion: 0.05
Nodes (63): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), bunTmp, inodeID (+55 more)

### Community 81 - "recordFlow"
Cohesion: 0.14
Nodes (16): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+8 more)

### Community 82 - "newTestTarget"
Cohesion: 0.27
Nodes (14): newTestTarget(), since(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded(), TestApplyFromWhitelistStopsFirst(), TestApplyOffStopsGuard(), TestApplyOffToWhitelistStartsGuard() (+6 more)

### Community 83 - "rcUser"
Cohesion: 0.12
Nodes (43): TestRemoveSSHAgentFacultativeKeepsGuardedSSHConfig(), rcTarget, BunTmpDir(), TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent() (+35 more)

### Community 84 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.13
Nodes (21): placedWalk, treeJudge, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedEntry(), placedFileStat() (+13 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 86 - "exe_supersede.h"
Cohesion: 0.22
Nodes (14): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+6 more)

### Community 87 - "loginServer"
Cohesion: 0.08
Nodes (23): markReady(), markReadyWithLogin(), notifySystemdReady(), hasSentinel(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel(), runtimePathsFor() (+15 more)

### Community 88 - "vaultfs_test.go"
Cohesion: 0.11
Nodes (22): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), IsBinaryContent(), IsOctalMode(), ParseMode(), inodeOf(), mustRead() (+14 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "runView"
Cohesion: 0.10
Nodes (5): daemonPage, runHost, runView, captureLogs(), orDefault()

### Community 92 - "FetchLatest"
Cohesion: 0.36
Nodes (9): installHooks(), releaseBuild(), Run(), staged(), FetchLatest(), fetchPair(), FetchVerified(), Release (+1 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 95 - "NetGuard"
Cohesion: 0.08
Nodes (14): foreignPinnedLink(), PumpRingbuf(), reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs() (+6 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "time.Duration"
Cohesion: 0.15
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 100 - "canonicalBinaryPath"
Cohesion: 0.15
Nodes (9): keptEntries, checkBtrfsKeys(), ensureBtrfsLayout(), canonicalBinaryPath(), diffBinaries(), TestDiffBinaries(), FirstOnBtrfs(), OnBtrfs() (+1 more)

### Community 101 - "root.go"
Cohesion: 0.09
Nodes (11): closeDumpFile(), Execute(), initLogger(), postRun(), stageItem, stageState, textAnswer, encodeTable() (+3 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "binaryList"
Cohesion: 0.11
Nodes (11): binaryList, watchTab, handCursors(), handCursorsOn(), showConfirm(), showCustomConfirm(), showDialog(), showErrorDialog() (+3 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "newTestVetter"
Cohesion: 0.08
Nodes (32): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+24 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "highlight_test.go"
Cohesion: 0.16
Nodes (17): fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens() (+9 more)

### Community 111 - "ConfinedEntry"
Cohesion: 0.13
Nodes (29): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+21 more)

### Community 112 - "installPage"
Cohesion: 0.07
Nodes (20): anyDaemonRunning(), actionGate, binSource, handButton, installAction, installedBinary, installPage, latestRelease (+12 more)

### Community 113 - "Vault"
Cohesion: 0.22
Nodes (4): Vault, Open(), Unopened(), WriteFileKeepMeta()

### Community 115 - "maintainerKeeper"
Cohesion: 0.18
Nodes (12): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp (+4 more)

### Community 116 - "deploy"
Cohesion: 0.18
Nodes (18): archAssets(), assetsFor(), checkUpdatePreconditions(), deploy(), downloadAndVerify(), fetchGUI(), guiAssetsFor(), hasAsset() (+10 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "UI"
Cohesion: 0.08
Nodes (22): nextConfig(), showChangelog(), channelOf(), confirmUpdate(), isTerminal(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI() (+14 more)

### Community 121 - "gui/session.go"
Cohesion: 0.15
Nodes (18): invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv() (+10 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "inet_connect.c"
Cohesion: 0.26
Nodes (15): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+7 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "classify"
Cohesion: 0.22
Nodes (12): call(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyOpenCallsExact(), arg(), attributesOf(), classify(), classifyCollection() (+4 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "catalogRefresher"
Cohesion: 0.12
Nodes (13): drain(), readEvents(), swapFile(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors() (+5 more)

### Community 131 - "logLine"
Cohesion: 0.13
Nodes (10): logHook, logLine, journalLine(), TestJournalLine(), logRow(), logStyle(), formatEntry(), sortedKeys() (+2 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "Finder"
Cohesion: 0.10
Nodes (29): findButton, span, findAll(), Finder, lowerRunes(), NewFinder(), newFindInput(), overlay() (+21 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "serveWebSocket"
Cohesion: 0.21
Nodes (8): configureWebSocket(), sameOrigin(), serveWebSocket(), TestSameOrigin(), TestSnapshotHubSingleViewerAndCoalescing(), writeSnapshot(), writeSnapshotLoop(), snapshotHub

### Community 143 - "serve.go"
Cohesion: 0.12
Nodes (22): basicAuth(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop(), relayServeSignals() (+14 more)

### Community 144 - "liveSession"
Cohesion: 0.26
Nodes (6): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), holdForward(), pingInterval(), liveSession

### Community 145 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.07
Nodes (45): class, lineStyles, pos, viewportState, classOf(), column(), dedent(), Edit() (+37 more)

### Community 146 - "Ledger"
Cohesion: 0.04
Nodes (50): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), bootstrapSelfProtectPlaceholders(), candidates(), confirmAll() (+42 more)

### Community 148 - "GuardEvent"
Cohesion: 0.05
Nodes (19): holdHeadless(), HeadlessLine(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary() (+11 more)

### Community 149 - "GlobReservations"
Cohesion: 0.15
Nodes (11): GlobChild, statFunc, compileGlobNames(), globText(), GlobReservations, resolveBits(), resolveChildren(), statKey() (+3 more)

### Community 150 - "find"
Cohesion: 0.22
Nodes (16): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), TestMovedFlagError(), movedFlagError(), Delete(), Find() (+8 more)

### Community 152 - "bufio.Reader"
Cohesion: 0.17
Nodes (14): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+6 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.05
Nodes (67): DrainUntilSignal(), CheckBPFLSM(), CheckEBPF(), backingDeviceUnion(), buildGuards(), buildOneGuard(), eventFilterOptions(), newPinGeneration() (+59 more)

### Community 155 - "watchSet"
Cohesion: 0.23
Nodes (6): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet

### Community 156 - "copytree.go"
Cohesion: 0.41
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 157 - "Daemon configuration"
Cohesion: 0.14
Nodes (13): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Launch rules: `[launch_rules]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Port sections, Process inspectors (+5 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 161 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 162 - "NetEventType"
Cohesion: 0.12
Nodes (11): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), newNetEventChecks(), FormatAddr(), NetEventTypes() (+3 more)

### Community 163 - "Gate"
Cohesion: 0.15
Nodes (6): placedKey(), Caller, Gate, Role, Listen(), Judge

### Community 164 - "debouncer"
Cohesion: 0.33
Nodes (4): TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), debouncer

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (48): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+40 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (56): fakeLink, homePage, pageDef, pageID, reel, resizeWatch, shell, snapEngines (+48 more)

### Community 167 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 168 - "procstats.go"
Cohesion: 0.38
Nodes (8): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic()

### Community 169 - "options.go"
Cohesion: 0.36
Nodes (9): InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks, StagedRelease (+1 more)

### Community 170 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "go_pkg_os"
Cohesion: 0.09
Nodes (29): auditEntry(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), groupWatchPaths(), relOrPath(), TestAuditTree() (+21 more)

### Community 175 - "TrustGuard"
Cohesion: 0.06
Nodes (23): tempGrant, trustRows, deleteResKeys(), setSupersedeTrust(), syncU32Map(), liveTrust(), cStr(), TestPlanTrustedDirs_DuplicateWriterAddsNoEmptyLoader() (+15 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 178 - "run"
Cohesion: 0.31
Nodes (9): deployGUI(), InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), syncGUI(), TestReadInstalledVersion() (+1 more)

### Community 179 - "runNetworkGuard"
Cohesion: 0.14
Nodes (18): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), WithSuspectsFrom() (+10 more)

### Community 180 - "AddServeFlags"
Cohesion: 0.36
Nodes (8): AddHeadlessFlag(), AddServeFlags(), TestParseServeFlagsRequiresCredentials(), init(), init(), init(), init(), init()

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 186 - "Binary trust"
Cohesion: 0.22
Nodes (9): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch rules, Launch scan, Multicall binaries, The dynamic loader (+1 more)

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "github.com/spf13/cobra.Command"
Cohesion: 0.14
Nodes (24): credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestSetupDumpLogRefusesSymlink() (+16 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 195 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 199 - "New"
Cohesion: 0.07
Nodes (51): harnessSuite, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD() (+43 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "pinstate_test.go"
Cohesion: 0.09
Nodes (33): configureDaemonLogging(), loadDaemonConfig(), noteReloadedSettings(), resolveConfigPath(), runLockdown(), pinOwnerLikelyAlive(), readPinState(), recoverPinState() (+25 more)

### Community 202 - "DaemonEvent"
Cohesion: 0.12
Nodes (28): UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemon(), runDaemonUI(), runHeadless() (+20 more)

### Community 203 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 204 - "EventType"
Cohesion: 0.05
Nodes (39): selfEventsFor(), snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting() (+31 more)

### Community 205 - "guardLSMHooks"
Cohesion: 0.25
Nodes (5): lsmHook, trustHook, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), TestSocketHooksAreRequired()

### Community 206 - "monitor"
Cohesion: 0.07
Nodes (30): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+22 more)

### Community 210 - "TestReadFrame"
Cohesion: 0.40
Nodes (4): TestReadFrame(), decodeFrame(), encodeWithSerial(), header()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "KernelDev"
Cohesion: 0.11
Nodes (27): fdStat, BinaryStat, physicalParent(), statFD(), binaryStatOf(), KernelDev(), StatBinary(), StatBinaryFile() (+19 more)

### Community 213 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "SecretStoreNotes"
Cohesion: 0.39
Nodes (7): dockerKeyringHelper(), gitCredentialHelpers(), readOwnedRegular(), SecretStoreNotes(), TestReadOwnedRegular_RefusesSymlinkAndOtherOwner(), TestSecretStoreNotes(), WarnSecretStores()

### Community 218 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 221 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (51): recordTickCmd(), recordModel, diffModel, editorModel, listenForDaemonEvents(), tickDaemonStats(), formatGuardEventLine(), sanitizeTerminalText() (+43 more)

### Community 228 - "mountTable"
Cohesion: 0.13
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), mountTable (+1 more)

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 238 - "Candidate"
Cohesion: 0.04
Nodes (69): appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories() (+61 more)

### Community 249 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 256 - "plan_test.go"
Cohesion: 0.07
Nodes (41): systemTool(), newLinkView(), checkLink(), linkFixture(), TestCheckLinksAccepts(), TestCheckLinksJudgesEachTarget(), TestCheckLinksRefusesUserFiles(), TestCheckLinksSeesThePlannedTree() (+33 more)

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 663 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **40 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `time.Duration`, `GuardInodeKey`, `canonicalBinaryPath`, `BinaryEntry`, `DaemonEvent`, `SanitizeText`, `TempBinary`, `TrustGuard`, `os.File`, `ConfinedEntry`, `GuardEvent`, `sync.Mutex`, `loginServer`, `buildOneGuard`, `time.Time`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.01733766233766234 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `parseDaemonEvents`, `go_pkg_strings`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09261790840738209 - nodes in this community are weakly interconnected._
- **Why does `monitor()` connect `monitor` to `time.Duration`, `secret_probe/main.go`, `How it works`, `Installation`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `op` be split into smaller, more focused modules?**
  _Cohesion score 0.11895161290322581 - nodes in this community are weakly interconnected._